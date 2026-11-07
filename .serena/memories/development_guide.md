# Development Guide

## Local run

Full stack (simplest):

```bash
docker compose up --build
```

Per-component (local dev, host):

```bash
# secret-assistant (gRPC crypto, :50052)
cd secret_assistant && go run ./cmd -port 50052

# redis (gateway needs it)
docker run -p 6379:6379 redis:7-bullseye

# gateway (HTTP API, :3000)
cd gateway && make run   # env: REDIS_HOST=localhost SECRET_ASSISTANT_URL=localhost:50052

# web (Vite dev, :5173; proxies /api -> http://localhost:3000, strips /api)
cd web && pnpm install && pnpm dev
```

`pnpm dev` relies on the Vite proxy (`/api` -> `http://localhost:3000`, `/api` prefix stripped), mirroring nginx, so the dev server reaches the gateway without the nginx container.

## Testing

Gateway (Go) — `cd gateway && make test` = `go test ./... -race -count=1 -cover`.

- `service_test.go`: hand-written fakes (`fakeCrypto`, `fakeStore`); one test per error path.
- `store_test.go`: `miniredis` (no real Redis); round-trip, consume-twice -> `ErrNotFound`, TTL applied, key namespacing.
- `handler_test.go`: package `handler_test` (black-box), `httptest` + real router via `server.New`; status codes, validation, error envelope, internal-error obfuscation.
- `config_test.go`: `t.Setenv` + `Load`; overrides and invalid values.
- `integration_test.go`: full stack over real TCP — httptest server + miniredis (TCP) + in-process fake gRPC; asserts the one-time guarantee (2nd read -> 404).
- Framework: `testify` (`assert`/`require`); no mock framework (hand-written fakes).

secret-assistant (Go): **no tests**; `go test ./...` / `go vet ./...` pass with "no test files".

Frontend (Vue/TS): **no tests**. Quality gates = `pnpm type-check` (vue-tsc), `pnpm build` (vite), `pnpm lint` (eslint + prettier).

## Lint / static analysis

- Gateway: `make lint` = `staticcheck ./...`; `make vet` = `go vet ./...`; `make fmt` = `gofmt`.
- secret-assistant: `go vet ./...` (added to CI; no staticcheck).
- web: `pnpm lint` (ESLint flat config + `@vue/eslint-config-typescript`, prettier formatting disabled via `skipFormatting`), `pnpm format` (prettier).

## Code style

- Gateway (Go): stdlib-first (`log/slog`, `errors` + `%w`, `errors.Is/As`); context propagated; handlers have 10s timeout (chi middleware). Sentinel errors (`store.ErrNotFound`, `service.ErrMessageNotFound`). Layering: handler (thin) -> service (logic) -> cryptoclient/store (adapters). `internal/messagepb` generated — never edit. `store` is a deep module (Send/Consume only).
- secret-assistant (Go): module path still `github.com/ReshetovItsMe/one-time-messaging-exchange-be`. The cross-service contract is an opaque `[]byte` blob (no JSON keys, no key material).
- web (TS): 4-space indent, prettier + eslint, pnpm (lockfile committed; CI uses `--frozen-lockfile`).

## CI

`.github/workflows/build.yml` (push to `main`, `chore/build-ci`):

- frontend: pnpm install (`--frozen-lockfile`) + `pnpm build` + `pnpm lint` (Node 26).
- gateway: `go mod verify` + `go build ./...` + `go vet ./...` + `go test ./... -race` (Go 1.27.x).
- secret_assistant: `go mod verify` + `go build` + `go test ./...` + `go vet ./...` (Go >=1.18).
- push_to_docker_hub: builds/pushes Docker images (needs Docker Hub secrets).

## Conventions / process

- OpenSpec (`openspec/`) drives spec changes; artifacts committed alongside code.
- Multi-agent team: `AGENTS-TEAM.md` (path ownership, contract ownership, branch/PR rules). One-off assistants ignore it.
- Conventional Commits, English, imperative mood.
