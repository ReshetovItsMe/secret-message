# Design: Secure key storage (envelope encryption, opaque blob)

## Context

The product's "complete confidentiality" claim is currently false: the RSA private key
is stored in Redis next to the ciphertext, so Redis access equals plaintext access. The
root cause is structural — the decryption key must be recoverable server-side (the link
carries no client key), and today that key travels through the gateway into Redis. The
fix must guarantee: **Redis compromise alone reveals no plaintext**, which requires a
secret that lives outside Redis and is held only by `secret_assistant`.

## Options considered

### (a) Opaque blob — crypto service returns an opaque blob; gateway knows nothing about keys

`Encrypt(text) → blob`, `Decrypt(blob) → text`. The gateway forwards and stores raw
bytes. **On its own this is not a security fix**: if the blob still embeds the private
key (as today's `EncryptedMessage` does), a Redis compromise still reveals plaintext.
It becomes a real fix only when combined with a service-held key (see below), at which
point it is the cleanest option because it also removes the JSON-in-proto-bytes seam and
the gateway's knowledge of key material.

### (b) Service key in env/secret — one RSA pair per service, private part only in secret-assistant

`Encrypt` wraps the ephemeral AES key with the service public key and returns only
`{ encryptedKey, data }`; the private key never leaves `secret_assistant`. This is the
correct security model and the smallest code delta (the existing RSA-OAEP + AES-GCM code
path is reused, only the key source changes). It keeps the JSON seam, however.

### (c) KMS / envelope encryption

Use an external KMS (AWS KMS, GCP KMS, Vault) to wrap/unwrap a data key per message.
This is the enterprise-correct answer: centralized key management, rotation, audit. It
is rejected here because it adds an external dependency, per-message network calls and
cloud lock-in for a small self-hosted (k3s/minikube) service. It remains the natural
upgrade path if the service grows or a key-management requirement appears.

### (d) privateKey encrypted at rest, stored separately from the ciphertext

Encrypt the per-message private key with a master key and store it under a different
key/column. This does not help against a full-Redis compromise (both pieces are in the
same store, same credentials) unless "separately" means a different storage system with
separate credentials — which this change does not introduce. Keeping per-message RSA
generation also adds no value once a service key exists. Rejected.

## Decision

**Option (a) with a service-held key — i.e. envelope encryption exposed as an opaque
blob.** It has the same security properties as (b) but additionally removes the
JSON-in-proto-bytes seam and the duplicated `EncryptedMessage`/`EncryptedPayload`
structs, which resolves the low-priority architectural question (JSON inside proto
bytes) in the same change.

- **Service key**: a single RSA-2048 keypair, private part PEM-encoded and supplied to
  `secret_assistant` only, via a file path selected by an env var (e.g.
  `RSA_PRIVATE_KEY_PATH`) so it works as a Docker secret / Kubernetes Secret mount. The
  public key is derived from the private key at startup. Keeping RSA-OAEP(SHA-256)
  preserves the scheme already documented in `AGENTS.md` ("RSA-OAEP + AES-GCM") and the
  existing code path.
- **Encrypt**: random 32-byte AES-256 key → AES-256-GCM encrypt plaintext (random nonce
  prepended, as today) → RSA-OAEP wrap the AES key with the service public key → emit a
  versioned blob `[version(1)][len(wrappedKey)(2)][wrappedKey][nonce][ciphertext+tag]`.
  No `privateKey`, no JSON.
- **Decrypt**: parse the blob (version check) → RSA-OAEP unwrap the AES key with the
  service private key → AES-GCM decrypt. Any length/version/auth failure is an error.
- **proto**: `Encrypt(EncryptRequest{plaintext}) → EncryptResponse{blob}`,
  `Decrypt(DecryptRequest{blob}) → DecryptResponse{plaintext}`; `blob` is `bytes`, no
  JSON. Both codegen outputs regenerated via `make proto` (gateway) and the
  `secret_assistant/README.md` command.
- **gateway**: `cryptoclient.Client` returns/takes `[]byte`; `store.Store` becomes
  `Send(ctx, blob []byte)` / `Consume(ctx, id) ([]byte, error)` and stores the blob
  directly as the Redis value (no JSON wrapper). The gateway never parses the blob.

## Threat model after the change

- Redis compromise alone → only wrapped data keys + ciphertext; no plaintext. **Fixed.**
- `secret_assistant` compromise → the private key is reachable; this is unavoidable
  (the service must hold the key). Mitigate with least-privilege secret access, a short
  message TTL, and (later) key rotation.
- Gateway compromise → gateway is already in the plaintext path (it calls `Encrypt`);
  unchanged.

## Non-goals / out of scope

- Key rotation, key versioning beyond a format version byte, multi-key support.
- KMS integration (recorded as a future option).
- Rewriting the delete-before-decrypt / at-most-once model (unchanged, per `AGENTS.md`).
- Transport encryption of gRPC/Redis (still in-cluster network only).

## Risks

- Contract change touches both services + both codegen outputs + docs; mitigate by
  doing it in one change and updating `AGENTS.md`'s API Contract section together.
- A malformed/legacy blob in Redis (messages written before deploy) will fail to
  decrypt after upgrade; acceptable — messages are one-time, TTL ≤ 24h, and this is a
  breaking change to an ephemeral store.
- Losing the service key makes all in-flight messages undecryptable; the key must be
  provisioned reliably (documented in the k8s change).
