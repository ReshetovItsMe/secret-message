// Command gateway is the API gateway of Secret Message: it exposes the HTTP
// API, talks to secret-assistant (gRPC) for crypto and Redis for storage.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ReshetovItsMe/secret-message/gateway/internal/config"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/cryptoclient"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/handler"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/logger"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/server"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/service"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/store"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gateway:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log := logger.New(cfg.LogPretty)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}
	defer rdb.Close()

	crypto, err := cryptoclient.New(cfg.SecretAssistantURL)
	if err != nil {
		return fmt.Errorf("init crypto client: %w", err)
	}
	defer crypto.Close()

	storage := store.New(rdb, cfg.MessageTTL, log)
	svc := service.New(crypto, storage)
	h := handler.New(svc, log)

	srv := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           server.New(h, log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("gateway listening", "addr", cfg.Addr())
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}
	case <-ctx.Done():
		log.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}
	}
	return nil
}
