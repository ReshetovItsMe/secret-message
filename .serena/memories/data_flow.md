# Data Flow

Lifecycle of a secret from creation to consumption across all four components.

## Write path (POST /message)

1. Browser -> `web` (nginx:8080) rewrites `/api/*` -> `/*` and proxies to `gateway:3000`.
2. `handler.Create` (gateway): decode JSON (`DisallowUnknownFields`, 1 MiB body cap), validate `message` (required, <=100000).
3. `service.CreateMessage`: `crypto.Encrypt(ctx, plaintext)` -> gRPC `Encrypt` to `secret-assistant`.
4. `secret-assistant` `Encrypt` (envelope): random 32-byte AES key -> AES-256-GCM encrypt plaintext (random nonce prepended to ciphertext) -> RSA-OAEP(SHA-256) wrap the AES key with the service public key -> versioned opaque blob `[version][len(wrappedKey)][wrappedKey][nonce][ciphertext+tag]` (no key material, no JSON).
5. `cryptoclient` returns the raw `[]byte` blob (opaque — the gateway never parses it).
6. `store.Send(ctx, blob)`: uuid -> key `msg:<uuid>`, `SET` the blob with TTL (default 24h). Returns id.
7. `handler.Create` -> `201 { messageId }`.
8. Frontend (`SecretInput` -> `SecretUrl`) shows link `<origin>/secret/<id>`.

NOTE: plaintext DOES transit the gateway (it is the Encrypt caller). Encryption is server-side (secret-assistant), NOT end-to-end. The service RSA private key stays in secret-assistant and is never stored in Redis; Redis holds only an opaque blob (wrapped key + ciphertext), so a Redis compromise alone reveals no plaintext. A secret-assistant (or its key) compromise, however, can reveal plaintext.

## Read path (GET /message?messageId=)

1. Browser -> web -> gateway `GET /message`.
2. `handler.Get`: require `messageId` query param (else 400).
3. `service.GetMessage`: `store.Consume(ctx, id)` FIRST.
   - `store.Consume` = Redis `GETDEL msg:<uuid>` (atomic read+delete = one-time guarantee). `redis.Nil` -> `store.ErrNotFound`.
4. Map `store.ErrNotFound` -> `service.ErrMessageNotFound` -> handler -> `404`.
5. On success, `crypto.Decrypt(ctx, blob)` -> gRPC `Decrypt` -> secret-assistant (RSA-OAEP unwrap AES key with the service private key -> AES-GCM decrypt; bad version/length/auth -> error).
6. `200 { message }`.

WARNING delete-before-decrypt: if `Decrypt` fails after `GETDEL` (secret-assistant outage, corrupt payload), the message is irrecoverably lost. Accepted tradeoff of the one-time guarantee (delete must be atomic with read). See `mem:architecture`.

## Frontend

- `/` (`MainView`): local reactive store toggles `SecretInput` (compose) / `SecretUrl` (show link).
- `/secret/:id` (`SecretView`): on mount `GET /message?messageId=:id`; 404 -> "already opened or does not exist"; loading/error states handled.
- axios `baseURL = '/api'`: dev via Vite proxy -> localhost:3000; prod via nginx rewrite.
- PWA: `/api/*` is `NetworkOnly` (never cached) to preserve one-time semantics.

## Storage

- One Redis key `msg:<uuid>` per message; TTL default 24h. No secondary indexes; no persistence config in compose.
- `store` is a deep module: only `Send`/`Consume`; key layout never leaks to callers.

## Error envelope

- `{ "error": { "code", "message" } }`; codes `BAD_REQUEST`, `VALIDATION_ERROR`, `NOT_FOUND`, `INTERNAL`.
- Client-facing messages generic; internal error strings never echoed (details only in logs).
