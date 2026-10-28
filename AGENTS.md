<!-- CODEGRAPH_START -->
## CodeGraph

In repositories indexed by CodeGraph (a `.codegraph/` directory exists at the repo root), reach for it BEFORE grep/find or reading files when you need to understand or locate code:

- **MCP tool** (when available): `codegraph_explore` answers most code questions in one call — the relevant symbols' verbatim source plus the call paths between them, including dynamic-dispatch hops grep can't follow. Name a file or symbol in the query to read its current line-numbered source. If it's listed but deferred, load it by name via tool search.
- **Shell** (always works): `codegraph explore "<symbol names or question>"` prints the same output.

If there is no `.codegraph/` directory, skip CodeGraph entirely — indexing is the user's decision.
<!-- CODEGRAPH_END -->

# Secret Message

One-time secret-sharing service: plaintext is encrypted (hybrid RSA + AES), a unique link is shared, and the recipient decrypts it exactly once — the stored payload is atomically deleted on read.

## Architecture

Four components wired together with Docker Compose, communicating over gRPC and HTTP:

```
Browser ──▶ web (Vue 3 + nginx, :8080)
              │  /api/* → gateway (nginx rewrite strips /api)
              ▼
           gateway (Go + chi, :3000)
              │  gRPC (insecure) ──▶ secret-assistant (Go, :50052)
              │                              │ encrypt/decrypt (RSA+AES)
              ▼                              ▼
           redis (:6379)              proto/message.proto (source of truth)
```

**Data flow:**

1. **POST /message** `{ "message": "text" }` — gateway calls gRPC `Encrypt` → secret-assistant returns a JSON `EncryptedMessage` (`privateKey` / `encryptedKey` / `data`, all `[]byte`), marshalled into the proto's `bytes body`. Gateway stores the whole payload under **one** Redis key `msg:<uuid>` with a TTL (default `24h`). Returns `201 { "messageId": "<uuid>" }`.
2. **GET /message?messageId=<uuid>** — gateway atomically `GETDEL msg:<uuid>` (read + delete in one round trip = the one-time guarantee), calls gRPC `Decrypt`, returns `200 { "message": "..." }`. A missing/consumed id returns `404`.
3. Frontend shows the secret at `/secret/:id` (SPA route; axios `baseURL` is `/api`).

## Tech Stack

| Component | Location | Stack | Port |
|---|---|---|---|
| `web` (frontend) | `web/` | Vue 3 (`<script setup lang="ts">`), Vite, vue-router, element-plus, @vueuse/core, axios, PWA, nginx SPA | 8080 |
| `gateway` (API) | `gateway/` | **Go 1.27**, chi, `log/slog`, go-redis v9, grpc-go, validator, samber/lo; testify + miniredis for tests | 3000 |
| `secret-assistant` (crypto) | `secret_assistant/` | Go 1.18, grpc-go, `crypto/aes` + `crypto/rsa` | 50052 |
| `redis` | — | redis:7-bullseye (no persistence config) | 6379 |

Note the Go version skew: gateway is `go 1.27`, secret-assistant is `go 1.18`.

## Directory Layout

- `proto/message.proto` — **source of truth** for the gRPC contract (`SecretAssistant.encrypt` / `.decrypt`). The `bytes body` fields carry a JSON-marshalled `EncryptedMessage`. Generated code is committed in **two places** with **two different codegen flows** (see Commands): `gateway/internal/messagepb/` and `secret_assistant/proto/`.
- `gateway/` — `cmd/gateway/main.go` (composition root), `internal/{config,server,handler,service,store,cryptoclient,messagepb,integration,logger}`. See `gateway/AGENTS.md` for module responsibilities.
- `secret_assistant/` — `cmd/main.go` (gRPC server, flag `-port`, default 50052), `internal/encrypt/`, `internal/decrypt/`, `pkg/messages/` (the `EncryptedMessage` struct — JSON keys are the cross-service contract), `proto/`.
- `web/src/` — `main.ts` (router: `/`, `/secret/:id`, catch-all redirect), `views/MainView.vue`, `views/SecretView.vue`, `components/` (SecretInput, SecretUrl, Header, Footer, icons).

## Common Commands

```bash
# Full stack
docker compose up --build

# Gateway (Go)
cd gateway && make build      # → ./bin/gateway
cd gateway && make run        # needs REDIS_HOST + SECRET_ASSISTANT_URL (see Env)
cd gateway && make test       # go test ./... -race -count=1 -cover
cd gateway && make lint       # staticcheck ./...
cd gateway && make proto      # regenerate internal/messagepb from ../proto/message.proto
cd gateway && make vet tidy

# secret-assistant (Go)
cd secret_assistant && go build ./... && go run ./cmd -port 50052
# regenerate its protos (see secret_assistant/README.md for the protoc invocation)

# Frontend (yarn)
cd web && yarn install && yarn dev              # vite dev server
cd web && yarn build && yarn type-check         # vue-tsc --noEmit + vite build
cd web && yarn lint
```

## Environment Variables

- **gateway (Go, `internal/config/config.go`):** `HOST` (0.0.0.0), `PORT` (3000), `REDIS_HOST` (redis), `REDIS_PORT` (6379), `REDIS_USERNAME`, `REDIS_PASSWORD`, `REDIS_DB` (0), `SECRET_ASSISTANT_URL` (secret-assistant:50052), `MESSAGE_TTL` (24h), `LOG_PRETTY` (false).
- **secret-assistant:** `-port` flag (50052).
- **web:** none (nginx proxies `/api/` → `gateway:3000`, strips the `/api` prefix).

## API Contract

- `POST /message` → `201 { "messageId" }`; body requires `message` (string, ≤100000 chars), unknown fields rejected.
- `GET /message?messageId=` → `200 { "message" }`; `404` for missing/consumed; `400` if `messageId` omitted.
- `GET /healthz` → `200 ok`.
- Error body: `{ "error": { "code", "message" } }` (e.g. `NOT_FOUND`, `INTERNAL`, `VALIDATION_ERROR`). Client-facing errors are generic; details only in logs.

## Conventions

- **Gateway (Go):** stdlib-first (`log/slog`, `errors` + `%w`, `errors.Is`/`As`). Context propagated everywhere; HTTP handlers have a 10s timeout. `store` is a deep module — keep its interface at `Send`/`Consume` and never leak the `msg:<uuid>` key layout to callers. `internal/messagepb` is generated; never edit by hand. New deps need a reason (chi, go-redis, grpc, validator, lo, testify/miniredis are the approved set).
- **secret-assistant (Go):** module path is still the old repo name `github.com/ReshetovItsMe/one-time-messaging-exchange-be`; the gateway uses `github.com/ReshetovItsMe/secret-message/gateway`. The JSON keys `privateKey`/`encryptedKey`/`data` are part of the cross-service contract — do not rename without a coordinated change on both sides.
- **web (TS):** 4-space indent, Prettier + ESLint per package; **yarn** (yarn.lock committed; CI uses `--frozen-lockfile`).

## Notes

- **OpenSpec** drives spec-driven changes in this repo: install it at the agent level (an OpenSpec plugin/skill for whatever agent you use — https://github.com/Fission-AI/OpenSpec), so the repo itself carries no agent-specific skill or command copies. OpenSpec writes its artifacts into `openspec/` and they are part of the codebase — commit them alongside the code they describe.
- gRPC is insecure (`insecure.NewCredentials()`); fine for the compose network, not for public exposure.
- CI runs on push to `main` and `chore/build-ci`; Node 18 is pinned for the `web` job, Go 1.27.x for the `gateway` job.
