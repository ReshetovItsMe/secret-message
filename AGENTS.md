<!-- CODEGRAPH_START -->
## CodeGraph

In repositories indexed by CodeGraph (a `.codegraph/` directory exists at the repo root), reach for it BEFORE grep/find or reading files when you need to understand or locate code:

- **MCP tool** (when available): `codegraph_explore` answers most code questions in one call — the relevant symbols' verbatim source plus the call paths between them, including dynamic-dispatch hops grep can't follow. Name a file or symbol in the query to read its current line-numbered source. If it's listed but deferred, load it by name via tool search.
- **Shell** (always works): `codegraph explore "<symbol names or question>"` prints the same output.

If there is no `.codegraph/` directory, skip CodeGraph entirely — indexing is the user's decision.
<!-- CODEGRAPH_END -->

# Secret Message

Secure one-time file/text sharing service: content is encrypted (hybrid RSA + AES), a unique link is shared, the recipient decrypts it once, and (by design) all encrypted data is deleted after access.

## Architecture

Four components wired together with Docker Compose, communicating over gRPC and HTTP:

```
Browser ──▶ web (Vue 3 + nginx, :8080)
              │  /api/* → gateway (nginx rewrite strips /api)
              ▼
           gateway (Hapi.js / Node TS, :3000)
              │  gRPC (insecure) ──▶ secret-assistant (Go, :50052)
              │                              │ encrypt/decrypt (RSA+AES)
              ▼                              ▼
           redis (:6379)              proto/message.proto (single source of truth)
```

**Data flow:**

1. **POST /message** `{ "message": "text" }` — gateway calls gRPC `Encrypt` → gets `{ encryptedKey, privateKey, data }` (JSON), then stores three linked Redis entries:
   - `<uuid>` → `encryptedKey` (the returned message id)
   - `encryptedKey` → `privateKey`
   - `privateKey` → `data`
   Returns `201 { messageId: "<uuid>" }`.
2. **GET /message?messageId=<uuid>** — gateway walks the Redis chain (id → encryptedKey → privateKey → data), calls gRPC `Decrypt`, returns `201 { message: "<decrypted text>" }`.
3. Frontend shows the secret at `/secret/:id` (SPA route; axios base URL is `/api`).

## Tech Stack

| Component | Location | Stack | Port |
|---|---|---|---|
| `web` (frontend) | `web/` | Vue 3 (`<script setup lang="ts">`), Vite, vue-router, element-plus, @vueuse/core, axios, PWA (register-service-worker), nginx SPA | 8080 |
| `gateway` (API) | `gateway/` | Node.js + TypeScript, Hapi.js, hapi-swagger (`/documentation`), @grpc/grpc-js, ioredis, joi, pino, uuid; jest | 3000 |
| `secret-assistant` (crypto) | `secret_assistant/` | Go 1.18, google.golang.org/grpc, crypto/aes + crypto/rsa | 50052 |
| `redis` | — | redis:7-bullseye (no persistence config, no TTLs) | 6379 |

## Directory Layout

- `proto/message.proto` — **source of truth** for the gRPC contract (`SecretAssistant.encrypt` / `.decrypt`). Generated code lives in `gateway/proto/` (JS, via `yarn build:proto`) and `secret_assistant/proto/` (Go, committed).
- `gateway/src/` — `index.ts` (Hapi bootstrap, CORS, swagger), `routes/message.ts` (POST/GET `/message`), `handlers/messages.ts` (orchestration + Redis chain), `services/db.ts` (ioredis), `services/secretAssistant.ts` (gRPC client, `SECRET_ASSISTANT_URL`).
- `secret_assistant/` — `cmd/main.go` (gRPC server, flag `-port`, default 50052), `internal/encrypt/encrypt.go` (hybrid RSA+AES: `generateAESKey`, `encryptAESKey`), `internal/decrypt/decrypt.go`, `pkg/messages/message.go` (payload types).
- `web/src/` — `main.ts` (router: `/`, `/secret/:id`, catch-all redirect), `views/MainView.vue` (input → URL), `views/SecretView.vue` (fetch + show secret), `components/` (SecretInput, SecretUrl, Header, Footer).

## Common Commands

```bash
# Full stack (dev compose builds local images)
docker compose up --build            # web:8080, gateway:3000, secret-assistant:50052, redis:6379

# Production (prebuilt images: reshetovitsme/secret-{frontend,gateway,assistant})
docker compose -f docker-compose.prod.yml up

# Gateway (yarn)
cd gateway && yarn install && yarn start        # dev (nodemon; expects REDIS_HOST + SECRET_ASSISTANT_URL)
cd gateway && yarn build                        # tsc → dist/
cd gateway && yarn build:proto                  # regenerate JS protos from ../proto/message.proto
cd gateway && yarn lint && yarn test            # eslint / jest

# Frontend (yarn)
cd web && yarn install && yarn dev              # vite dev server
cd web && yarn build && yarn type-check         # vue-tsc --noEmit + vite build
cd web && yarn lint

# Go service
cd secret_assistant && go build ./... && go run ./cmd -port 50052
```

## Environment Variables

- **gateway:** `HOST` (0.0.0.0), `PORT` (3000), `REDIS_HOST` (redis), `REDIS_PORT` (6379), `REDIS_USERNAME`, `REDIS_PASSWORD`, `REDIS_DB` (0), `SECRET_ASSISTANT_URL` (secret-assistant:50052)
- **secret-assistant:** `-port` flag (50052)
- **web:** none (nginx proxies `/api/` → `gateway:3000`)

## CI/CD

`.github/workflows/build.yml` runs on push to `main` (and `chore/build-ci`): web (install → build → lint), gateway (install → build → lint → test), secret-assistant (Go build/test). Node 18 in CI.

## Conventions

- TypeScript: 4-space indent, Prettier + ESLint configured per package; **yarn** as package manager (yarn.lock committed; CI uses `--frozen-lockfile`).
- Go: standard layout (`cmd/`, `internal/`, `pkg/`), `go 1.18`, grpc reflection enabled in dev.
- API responses: both routes return HTTP `201` (including GET).
- Swagger UI available at `/documentation` on the gateway.

## Known Gaps & Notes (verify before relying on them)

- **One-time deletion is NOT implemented.** README claims data is deleted after access, but there is no `DEL` and no TTL anywhere in `gateway/src` — Redis keys persist forever. Implement expiry/deletion if the one-time guarantee matters.
- **No tests:** gateway `jest` runs with `--passWithNoTests`; Go service has no test files. CI "testing" steps are effectively no-ops.
- **Validation is permissive:** `POST /message` Joi schema is an empty `Joi.object()` — any payload is accepted; handler reads `payload.message` unchecked.
- **Module name drift:** `secret_assistant/go.mod` declares `github.com/ReshetovItsMe/one-time-messaging-exchange-be` (old repo name); the repo is now `secret-message`.
- **Insecure gRPC:** `credentials.createInsecure()` — fine for the compose network, not for public exposure.
- Generated proto code is committed in both `gateway/proto/` and `secret_assistant/proto/`; regenerate after editing `proto/message.proto`.
