# Domain Model & Glossary

## Business entities

- **Secret (Message)** — the plaintext a sender wants to share exactly once. Plain string, <= 100000 chars. Never stored in plaintext: only its encrypted form exists in Redis.
- **MessageID** — UUID (v4) identifying a stored secret. Returned by `POST /message` as `messageId`; consumed by `GET /message?messageId=`. Maps to Redis key `msg:<uuid>`.
- **Encrypted blob** — the sealed, opaque, versioned binary representation stored in Redis and passed over gRPC. Produced by `secret_assistant` `Encrypt` (wrapped AES-256-GCM data key + ciphertext); the gateway forwards and stores it as raw `[]byte` without parsing. Contains no private key material.
- **Store** — deep storage module; interface `Send(ctx, blob []byte) (id, error)` / `Consume(ctx, id) ([]byte, error)`. Redis-backed; `Consume` = `GETDEL` (atomic read+delete). Returns `ErrNotFound` when missing or already consumed.
- **Service** — business logic; `CreateMessage(ctx, plaintext) (id, error)` (encrypt -> store) and `GetMessage(ctx, id) (plaintext, error)` (consume -> decrypt). Maps `store.ErrNotFound` -> `ErrMessageNotFound`.
- **CryptoClient** — gateway adapter at the crypto seam; `Encrypt(ctx, plaintext) ([]byte, error)` / `Decrypt(ctx, blob []byte) (plaintext, error)` over gRPC to secret-assistant. The blob is opaque to the gateway.

## HTTP API

- `POST /message` `{message}` -> `201 {messageId}`; unknown fields rejected; `message` required, <= 100000.
- `GET /message?messageId=` -> `200 {message}`; `404` missing/consumed; `400` if `messageId` omitted.
- `GET /healthz` -> `200 ok`.
- Error body `{error:{code,message}}`; codes `NOT_FOUND`, `INTERNAL`, `VALIDATION_ERROR`, `BAD_REQUEST`. Client-facing errors are generic; details only in logs.

## Glossary

- **One-time / one-time guarantee** — a message can be read exactly once; reading atomically deletes it.
- **GETDEL** — Redis command returning a key's value and deleting it in one atomic operation; the mechanism behind the one-time guarantee.
- **TTL (Time To Live)** — how long an unread message survives before Redis evicts it; default 24h (`MESSAGE_TTL`).
- **Envelope encryption** — a per-message AES-256-GCM key encrypts the plaintext; a service-held RSA-2048 key wraps that AES key.
- **AES** — symmetric cipher encrypting the plaintext body.
- **RSA** — asymmetric cipher; the service-held RSA-2048 key wraps the per-message AES data key.
- **gRPC** — RPC framework; gateway <-> secret-assistant transport (insecure, in-cluster).
- **SecretAssistant** — the gRPC service (`encrypt`/`decrypt`) implemented by secret-assistant; contract in `proto/message.proto`.
- **messageId** — the UUID a sender shares with the recipient (the "secret link" id).
- **Consume vs Read** — `Consume` deletes on read (`GETDEL`); a plain "read" is deliberately not exposed because it would break the one-time property.
