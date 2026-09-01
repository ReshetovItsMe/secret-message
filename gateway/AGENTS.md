# Gateway (Go) — Agent Context

Working notes for AI agents and humans editing this service.

## What this is

HTTP gateway of Secret Message (Go 1.27, chi router, slog, go-redis, grpc-go).
Serves `POST /message` / `GET /message?messageId=` and talks to
secret-assistant (gRPC) + Redis. See [README.md](README.md) for the full
picture.

## Layout

- `cmd/gateway/main.go` — composition root: config → logger → redis → gRPC → store → service → handler → server, graceful shutdown
- `internal/handler` — thin HTTP layer: decode, validate, map errors to status codes. No business logic here.
- `internal/service` — the deep module: `CreateMessage` (encrypt → store), `GetMessage` (consume → decrypt). Sentinel: `service.ErrMessageNotFound`.
- `internal/store` — `MessageStore` interface + Redis implementation. **Deep module: keep the interface at `Send`/`Consume`.** Single key per message (`msg:<uuid>`), TTL, atomic `GETDEL` (one-time guarantee). Do not leak the key layout to callers.
- `internal/cryptoclient` — gRPC adapter. `EncryptedPayload` is the shared shape; JSON keys (`privateKey`/`encryptedKey`/`data`) are part of the contract with secret-assistant — do not rename without a coordinated change on the Go crypto service.
- `internal/messagepb` — generated protobuf. Never edit; regenerate with `make proto`.
- `internal/integration` — end-to-end test over real TCP (httptest + miniredis + in-process fake gRPC service).

## Conventions

- Standard library first: `log/slog` for logging, `errors` with `%w` wrapping, `errors.Is`/`errors.As`.
- `github.com/samber/lo` is available for slice/map helpers and `lo.Must` for infallible operations; use it when it removes a loop, not for its own sake.
- Validation with `go-playground/validator` on request DTOs; unknown JSON fields are rejected (`DisallowUnknownFields`).
- Errors to clients are generic (`INTERNAL`); details only in logs. Never echo internal error strings.
- Context is propagated everywhere; HTTP handlers have a 10s timeout via middleware.
- New dependencies need a reason; prefer stdlib + the existing small set (chi, go-redis, grpc, validator, lo, testify/miniredis in tests).

## Commands

```bash
make build   # ./bin/gateway
make test    # go test ./... -race
make lint    # staticcheck ./...
make proto   # regenerate internal/messagepb from ../proto/message.proto
```

## Testing

- Handler tests hit the real router via `httptest` with stubbed service adapters.
- Store tests use `miniredis` (no real Redis needed).
- Service tests use hand-written fakes (no mock framework).
- The integration test exercises the full stack over TCP and asserts the
  one-time guarantee (second read → 404).
- Keep the one-time semantics covered: consume twice → `ErrNotFound`.
