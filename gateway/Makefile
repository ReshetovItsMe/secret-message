# Secret Message — gateway (Go)
#
# Common targets:
#   make build   — compile the binary into ./bin/gateway
#   make run     — run locally (needs REDIS_HOST/SECRET_ASSISTANT_URL, see README)
#   make test    — run all tests with the race detector
#   make lint    — run golangci-lint
#   make lint-staticcheck — run standalone staticcheck
#   make proto   — regenerate protobuf code from ../proto/message.proto
#   make tidy    — go mod tidy
#   make vet     — go vet

GO       ?= go
BINARY    = bin/gateway
LINTER   ?= golangci-lint

.PHONY: build run test lint lint-staticcheck proto tidy vet fmt clean

build:
	$(GO) build -o $(BINARY) ./cmd/gateway

run: build
	./$(BINARY)

test:
	$(GO) test ./... -race -count=1 -cover

lint:
	$(LINTER) run

lint-staticcheck:
	staticcheck ./...

proto:
	protoc -I ../proto \
		--go_out=. \
		--go_opt=module=github.com/ReshetovItsMe/secret-message/gateway \
		--go_opt=Mmessage.proto=github.com/ReshetovItsMe/secret-message/gateway/internal/messagepb \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/ReshetovItsMe/secret-message/gateway \
		--go-grpc_opt=Mmessage.proto=github.com/ReshetovItsMe/secret-message/gateway/internal/messagepb \
		../proto/message.proto

tidy:
	$(GO) mod tidy

vet:
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -rf $(BINARY)
