# Gateway (Go)

HTTP API gateway of Secret Message. Written in Go 1.27 using the `chi` router.

It is the orchestration layer: it exposes the HTTP API, calls the
`secret-assistant` gRPC service for encryption/decryption, and stores the
sealed payload in Redis.

## Endpoints

- `POST /message` `{ "message": "text" }` → `201 { "messageId": "<uuid>" }`
- `GET /message?messageId=<uuid>` → `200 { "message": "..." }` (one-time: the
  payload is atomically deleted on read); `404` if missing/consumed
- `GET /healthz` → `200 ok`

See `AGENTS.md` (repo root) for the full architecture and `AGENTS.md` (this
directory) for module responsibilities and conventions.

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `HOST` | `0.0.0.0` | Listen address |
| `PORT` | `3000` | Listen port |
| `REDIS_HOST` | `redis` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_USERNAME` | — | Redis username (optional) |
| `REDIS_PASSWORD` | — | Redis password (optional) |
| `REDIS_DB` | `0` | Redis database |
| `SECRET_ASSISTANT_URL` | `secret-assistant:50052` | gRPC crypto service `host:port` |
| `MESSAGE_TTL` | `24h` | TTL of a stored message before Redis evicts it |
| `LOG_PRETTY` | `false` | Human-readable logs (dev) instead of JSON |

## Commands

```bash
make build   # compile → ./bin/gateway
make run     # build + run (needs REDIS_HOST + SECRET_ASSISTANT_URL)
make test    # go test ./... -race -count=1 -cover
make lint    # staticcheck ./...
make proto   # regenerate internal/messagepb from ../proto/message.proto
make vet     # go vet ./...
make tidy    # go mod tidy
```
