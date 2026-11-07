# Proposal: Secure key storage (envelope encryption, opaque blob)

## Why

Issue #15. Today `secret_assistant` generates a fresh RSA-2048 keypair per message and
returns the **private key** as part of `EncryptedMessage`, which the gateway stores in
Redis right next to the ciphertext (`privateKey` + `encryptedKey` + `data`). Whoever
can read Redis can therefore recover every plaintext ever stored, which contradicts the
README's promise of "complete confidentiality". The decryption key must move out of
Redis and into the crypto service, so that a Redis compromise alone no longer exposes
plaintext.

## What Changes

- **BREAKING** Replace the per-message RSA keypair with a single service-held RSA key:
  `Encrypt` wraps an ephemeral AES-256 key with the service *public* key, `Decrypt`
  unwraps it with the service *private* key. The private key never leaves
  `secret_assistant`.
- **BREAKING** `Encrypt` returns an opaque, versioned binary blob instead of the JSON
  `{ privateKey, encryptedKey, data }`. The `privateKey`/`encryptedKey`/`data` JSON
  contract is removed.
- **BREAKING** `proto/message.proto` carries an opaque `bytes blob` (no JSON inside
  bytes); message fields renamed for clarity. Both codegen outputs regenerated.
- Gateway `cryptoclient` and `store` become opaque: they forward/store raw `[]byte`
  blobs and have no knowledge of keys or ciphertext structure. This also removes the
  duplicated `EncryptedMessage`/`EncryptedPayload` structs across the service seam.
- The service private key is supplied to `secret_assistant` only (env-selected PEM file
  / mounted secret); never committed, logged, or stored in Redis.
- Update `AGENTS.md` (crypto description + contract), `README.md` (honest security
  claim), and `.serena/memories/*` to match.

## Capabilities

### New Capabilities

- `secure-key-storage`: stored secret payloads contain no private key material, and a
  Redis compromise alone does not reveal plaintext.

## Impact

- `secret_assistant/**` (encrypt, decrypt, cmd, config for the service key).
- `gateway/internal/{cryptoclient,store,service,cmd}` (opaque blob plumbing).
- `proto/message.proto` + regenerated `gateway/internal/messagepb/**` and
  `secret_assistant/proto/**`.
- `AGENTS.md`, `README.md`, `docker-compose*.yml` (key mount), `.serena/memories/*`.

## Sequencing

Depends on `crypto-tests` (regression tests land first and are adjusted to the new
signature). Service-key provisioning lives in this change's own docker-compose task
group (5); there is no Kubernetes Secret, since deployment/CD is out of scope for the
project. Independent of `ci-pipeline`.
