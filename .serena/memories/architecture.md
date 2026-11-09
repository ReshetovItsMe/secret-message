# Architecture

One-time secret-sharing service: a sender encrypts plaintext (envelope encryption — a per-message AES-256-GCM key wrapped by a service-held RSA-2048 key), gets a unique link; the recipient opens it and the plaintext is decrypted exactly once — the stored payload is atomically deleted on read.

## Components

| Component | Location | Stack | Port | Role |
|---|---|---|---|---|
| web | `web/` | Vue 3 (`<script setup lang="ts">`), Vite, vue-router, element-plus, @vueuse/core, axios, PWA, nginx SPA | 8080 | Sender/recipient UI |
| gateway | `gateway/` | Go 1.27, chi, log/slog, go-redis v9, grpc-go, validator, samber/lo | 3000 | HTTP API + orchestration |
| secret-assistant | `secret_assistant/` | Go 1.27.0, grpc-go, crypto/aes + crypto/rsa | 50052 | Crypto (encrypt/decrypt) gRPC |
| redis | — | redis:7-bullseye | 6379 | One-time payload store (TTL) |

Wired with Docker Compose. Browser -> web (nginx) -> gateway (chi; `/api/*` rewritten, `/api` prefix stripped) -> gRPC -> secret-assistant; gateway -> Redis.

## Request flow

1. `POST /message { message }` — gateway calls gRPC `Encrypt`; secret-assistant returns an opaque, versioned binary blob (wrapped AES-256-GCM data key + ciphertext; no key material, no JSON). Gateway stores the blob under one Redis key `msg:<uuid>` with TTL (default 24h). Returns `201 { messageId }`.
2. `GET /message?messageId=<uuid>` — gateway atomically `GETDEL msg:<uuid>` (read+delete in one round trip = the one-time guarantee), then gRPC `Decrypt`, returns `200 { message }`. Missing/consumed -> `404`.
3. Frontend shows the secret at `/secret/:id` (SPA route; axios `baseURL` = `/api`).

## Key invariants & concepts

- One-time guarantee = `GETDEL` (atomic read+delete). Consume happens BEFORE decrypt; if decrypt fails after consume, the message is irrecoverable — accepted tradeoff of delete-on-read.
- Crypto is envelope encryption, done server-side in secret-assistant: a per-message AES-256-GCM key encrypts the plaintext, and a single service-held RSA-2048/OAEP key wraps that AES key. The service private key never leaves secret-assistant and is never stored in Redis — Redis holds only an opaque blob, so a Redis compromise alone reveals no plaintext. Plaintext transits the gateway (it calls `Encrypt`) — this is NOT end-to-end encryption.
- `store` is a deep module: callers only see `Send`/`Consume`; the `msg:<uuid>` key layout never leaks.
- gRPC is insecure (`insecure.NewCredentials()`) — compose network only, not public exposure.
- Both Go modules target Go 1.27.0; CI uses Go 1.27.x.

## Cross-service contracts

- `proto/message.proto` — source of truth for the gRPC contract (`SecretAssistant.encrypt`/`.decrypt`). Requests/responses carry an opaque `bytes blob` (no JSON inside bytes): `EncryptRequest{plaintext}` → `EncryptResponse{blob}`, `DecryptRequest{blob}` → `DecryptResponse{plaintext}`. Generated code is committed in TWO places with TWO codegen flows: `gateway/internal/messagepb/` and `secret_assistant/proto/`.
- The cross-service contract is an opaque `[]byte` blob: the gateway forwards and stores it without parsing, and no key or ciphertext field crosses the seam. There is no `EncryptedMessage`/`EncryptedPayload` struct on either side. Contract changes must be coordinated across both services and both codegen outputs.
- Module paths: gateway `github.com/ReshetovItsMe/secret-message/gateway`; secret-assistant still `github.com/ReshetovItsMe/one-time-messaging-exchange-be`.

## Directory layout (gateway)

`gateway/` — `cmd/gateway/main.go` (composition root), `internal/{config,server,handler,service,store,cryptoclient,messagepb,integration,logger}`. Layering: handler (thin HTTP) -> service (business logic) -> cryptoclient (gRPC adapter) / store (Redis adapter). `internal/messagepb` is generated — never edit by hand.

## Environment

- gateway: `HOST`(0.0.0.0), `PORT`(3000), `REDIS_HOST`(redis), `REDIS_PORT`(6379), `REDIS_USERNAME`, `REDIS_PASSWORD`, `REDIS_DB`(0), `SECRET_ASSISTANT_URL`(secret-assistant:50052), `MESSAGE_TTL`(24h), `LOG_PRETTY`(false).
- secret-assistant: `-port` flag (50052); service RSA private key via `RSA_PRIVATE_KEY_PATH` (PEM file, never committed).
- web: none in prod (nginx proxies `/api/` -> `gateway:3000`, strips `/api`); in dev, Vite proxies `/api` -> `http://localhost:3000` (same rewrite).

## Commands

- Full stack: `docker compose up --build`
- gateway: `make build` / `make run` / `make test` / `make lint` / `make proto` / `make vet tidy` (see `gateway/AGENTS.md`)
- secret-assistant: `go build ./... && go run ./cmd -port 50052`
- web (pnpm): `pnpm install && pnpm dev` / `pnpm build && pnpm type-check` / `pnpm lint`

## Related

Business entities and glossary: `mem:domain_model`.
Process/team conventions: `AGENTS.md`, `AGENTS-TEAM.md`, `gateway/AGENTS.md`.
Spec-driven changes via OpenSpec in `openspec/`.
