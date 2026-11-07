# Architecture

One-time secret-sharing service: a sender encrypts plaintext (hybrid RSA+AES), gets a unique link; the recipient opens it and the plaintext is decrypted exactly once — the stored payload is atomically deleted on read.

## Components

| Component | Location | Stack | Port | Role |
|---|---|---|---|---|
| web | `web/` | Vue 3 (`<script setup lang="ts">`), Vite, vue-router, element-plus, @vueuse/core, axios, PWA, nginx SPA | 8080 | Sender/recipient UI |
| gateway | `gateway/` | Go 1.27, chi, log/slog, go-redis v9, grpc-go, validator, samber/lo | 3000 | HTTP API + orchestration |
| secret-assistant | `secret_assistant/` | Go 1.18, grpc-go, crypto/aes + crypto/rsa | 50052 | Crypto (encrypt/decrypt) gRPC |
| redis | — | redis:7-bullseye | 6379 | One-time payload store (TTL) |

Wired with Docker Compose. Browser -> web (nginx) -> gateway (chi; `/api/*` rewritten, `/api` prefix stripped) -> gRPC -> secret-assistant; gateway -> Redis.

## Request flow

1. `POST /message { message }` — gateway calls gRPC `Encrypt`; secret-assistant returns a JSON `EncryptedMessage` (`privateKey`/`encryptedKey`/`data`, all `[]byte`) marshalled into the proto `bytes body`. Gateway stores the whole payload under one Redis key `msg:<uuid>` with TTL (default 24h). Returns `201 { messageId }`.
2. `GET /message?messageId=<uuid>` — gateway atomically `GETDEL msg:<uuid>` (read+delete in one round trip = the one-time guarantee), then gRPC `Decrypt`, returns `200 { message }`. Missing/consumed -> `404`.
3. Frontend shows the secret at `/secret/:id` (SPA route; axios `baseURL` = `/api`).

## Key invariants & concepts

- One-time guarantee = `GETDEL` (atomic read+delete). Consume happens BEFORE decrypt; if decrypt fails after consume, the message is irrecoverable — accepted tradeoff of delete-on-read.
- Crypto is hybrid RSA+AES, done server-side in secret-assistant: RSA wraps a per-message AES key; AES-GCM encrypts the plaintext. Plaintext transits the gateway (it calls `Encrypt`) and never reaches Redis — this is NOT end-to-end encryption; the RSA private key is stored alongside the ciphertext.
- `store` is a deep module: callers only see `Send`/`Consume`; the `msg:<uuid>` key layout never leaks.
- gRPC is insecure (`insecure.NewCredentials()`) — compose network only, not public exposure.
- Go version skew: gateway `go 1.27`, secret-assistant `go 1.18`.

## Cross-service contracts

- `proto/message.proto` — source of truth for the gRPC contract (`SecretAssistant.encrypt`/`.decrypt`). `bytes body` carries a JSON-marshalled `EncryptedMessage`. Generated code is committed in TWO places with TWO codegen flows: `gateway/internal/messagepb/` and `secret_assistant/proto/`.
- JSON keys `privateKey`/`encryptedKey`/`data` are the cross-service contract. Struct lives in `secret_assistant/pkg/messages/` (`EncryptedMessage`), mirrored in `gateway/internal/cryptoclient/` (`EncryptedPayload`). Do not rename without a coordinated change on both sides.
- Module paths: gateway `github.com/ReshetovItsMe/secret-message/gateway`; secret-assistant still `github.com/ReshetovItsMe/one-time-messaging-exchange-be`.

## Directory layout (gateway)

`gateway/` — `cmd/gateway/main.go` (composition root), `internal/{config,server,handler,service,store,cryptoclient,messagepb,integration,logger}`. Layering: handler (thin HTTP) -> service (business logic) -> cryptoclient (gRPC adapter) / store (Redis adapter). `internal/messagepb` is generated — never edit by hand.

## Environment

- gateway: `HOST`(0.0.0.0), `PORT`(3000), `REDIS_HOST`(redis), `REDIS_PORT`(6379), `REDIS_USERNAME`, `REDIS_PASSWORD`, `REDIS_DB`(0), `SECRET_ASSISTANT_URL`(secret-assistant:50052), `MESSAGE_TTL`(24h), `LOG_PRETTY`(false).
- secret-assistant: `-port` flag (50052).
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
