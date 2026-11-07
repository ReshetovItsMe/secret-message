# Domain Model & Glossary

## Business entities

- **Secret (Message)** — the plaintext a sender wants to share exactly once. Plain string, <= 100000 chars. Never stored in plaintext: only its encrypted form exists in Redis.
- **MessageID** — UUID (v4) identifying a stored secret. Returned by `POST /message` as `messageId`; consumed by `GET /message?messageId=`. Maps to Redis key `msg:<uuid>`.
- **EncryptedMessage / EncryptedPayload** — the sealed representation stored in Redis and passed over gRPC. JSON keys (cross-service contract): `privateKey` (RSA private key, `[]byte`), `encryptedKey` (AES key wrapped by RSA, `[]byte`), `data` (ciphertext, `[]byte`). Struct: `secret_assistant/pkg/messages/message.go` (`EncryptedMessage`); mirror in gateway `internal/cryptoclient` (`EncryptedPayload`).
- **Store** — deep storage module; interface `Send(ctx, payload) (id, error)` / `Consume(ctx, id) (payload, error)`. Redis-backed; `Consume` = `GETDEL` (atomic read+delete). Returns `ErrNotFound` when missing or already consumed.
- **Service** — business logic; `CreateMessage(ctx, plaintext) (id, error)` (encrypt -> store) and `GetMessage(ctx, id) (plaintext, error)` (consume -> decrypt). Maps `store.ErrNotFound` -> `ErrMessageNotFound`.
- **CryptoClient** — gateway adapter at the crypto seam; `Encrypt(ctx, plaintext) (payload, error)` / `Decrypt(ctx, payload) (plaintext, error)` over gRPC to secret-assistant.

## HTTP API

- `POST /message` `{message}` -> `201 {messageId}`; unknown fields rejected; `message` required, <= 100000.
- `GET /message?messageId=` -> `200 {message}`; `404` missing/consumed; `400` if `messageId` omitted.
- `GET /healthz` -> `200 ok`.
- Error body `{error:{code,message}}`; codes `NOT_FOUND`, `INTERNAL`, `VALIDATION_ERROR`, `BAD_REQUEST`. Client-facing errors are generic; details only in logs.

## Glossary

- **One-time / one-time guarantee** — a message can be read exactly once; reading atomically deletes it.
- **GETDEL** — Redis command returning a key's value and deleting it in one atomic operation; the mechanism behind the one-time guarantee.
- **TTL (Time To Live)** — how long an unread message survives before Redis evicts it; default 24h (`MESSAGE_TTL`).
- **Hybrid encryption** — RSA encrypts a per-message symmetric AES key; AES encrypts the actual plaintext.
- **AES** — symmetric cipher encrypting the plaintext body.
- **RSA** — asymmetric cipher wrapping the per-message AES key.
- **gRPC** — RPC framework; gateway <-> secret-assistant transport (insecure, in-cluster).
- **SecretAssistant** — the gRPC service (`encrypt`/`decrypt`) implemented by secret-assistant; contract in `proto/message.proto`.
- **messageId** — the UUID a sender shares with the recipient (the "secret link" id).
- **Consume vs Read** — `Consume` deletes on read (`GETDEL`); a plain "read" is deliberately not exposed because it would break the one-time property.
