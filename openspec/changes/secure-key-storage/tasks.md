# Tasks: Secure key storage (envelope encryption, opaque blob)

Role ownership: `backend` = `secret_assistant/**`, `gateway/**`, `proto/**`;
`devops` = `docker-compose*.yml`, `**/Dockerfile`, `.github/**`, root `Makefile`;
`teamlead` = decisions + `AGENTS.md`/`README.md`/`.serena/memories/*`.

## 1. Crypto service: envelope encryption (backend)

- [x] 1.1 Load a single service RSA-2048 private key from a PEM file selected by an env
  var (e.g. `RSA_PRIVATE_KEY_PATH`), deriving the public key at startup. Fail fast if
  absent or invalid.
- [x] 1.2 Rewrite `Encrypt` to wrap an ephemeral AES-256 key with the service public key
  and emit a versioned opaque blob (no `privateKey`, no JSON).
- [x] 1.3 Rewrite `Decrypt` to parse the blob, unwrap the AES key with the service
  private key and AES-GCM-decrypt; reject bad version/length/auth with an error.
- [x] 1.4 Remove the `privateKey` field and JSON round-trip from
  `secret_assistant/pkg/messages` and `cmd/main.go`.

## 2. Contract: proto + codegen (backend)

- [x] 2.1 Update `proto/message.proto` to opaque `bytes blob` fields with clear message
  names (`EncryptRequest{plaintext}`, `EncryptResponse{blob}`,
  `DecryptRequest{blob}`, `DecryptResponse{plaintext}`).
- [x] 2.2 Regenerate both codegen outputs (`make proto` in gateway and the
  `secret_assistant/README.md` command); verify both sets of `.pb.go` files.

## 3. Gateway: opaque plumbing (backend)

- [ ] 3.1 Change `cryptoclient.Client` to accept/return `[]byte` blobs and drop the
  `EncryptedPayload` struct and its JSON (un)marshal.
- [ ] 3.2 Change `store.Store` to `Send(ctx, blob []byte)` / `Consume(ctx, id) ([]byte,
  error)` and store the blob directly as the Redis value (no JSON wrapper).
- [ ] 3.3 Update `service`, `cmd/gateway/main.go`, and the integration fake to the new
  opaque signatures; keep the HTTP surface (`POST/GET /message`, error envelope)
  unchanged.
- [ ] 3.4 `cd gateway && make build && make test && make lint && make vet` green.

## 4. Tests (backend, extends crypto-tests)

- [ ] 4.1 `secret_assistant`: update/add tests for the new signature — round-trip,
  wrong service key → error, tampered blob → error, empty and large inputs.
- [ ] 4.2 `cd secret_assistant && go build ./... && go test ./... && go vet ./...`
  green.

## 5. Key provisioning (devops)

- [x] 5.1 Add the service-key mount to `docker-compose.yml` /
  `docker-compose.prod.yml` and document how a developer supplies a synthetic key; do
  not commit any real key material.
- [x] 5.2 `docker compose config --quiet` valid.

## 6. Documentation and contract (teamlead)

- [ ] 6.1 Update `AGENTS.md` API Contract + crypto description (envelope encryption,
  opaque blob, remove the "private key stored alongside ciphertext" note).
- [ ] 6.2 Update `README.md` to an honest confidentiality statement.
- [ ] 6.3 Update `.serena/memories/{architecture,data_flow}.md` to match.
