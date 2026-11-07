# Proposal: Crypto module tests

## Why

Issue #31. `secret_assistant` ships zero tests: `go test ./...` reports "[no test files]".
`internal/encrypt` and `internal/decrypt` are pure functions (hybrid RSA-OAEP + AES-GCM,
no I/O or external dependencies), so they are trivially unit-testable, yet this is the
security-critical heart of the product and a regression in `Encrypt`/`Decrypt` would ship
silently today. Before `secure-key-storage` rewrites the signature (per-message RSA key →
service-held key, opaque blob), we pin the current behaviour with tests so that refactor
is provably non-regressive.

## What Changes

- Add table-driven unit tests for `secret_assistant/internal/encrypt` and
  `secret_assistant/internal/decrypt`, against the current signature
  (`Encrypt(text) → *EncryptedMessage`, `Decrypt(*EncryptedMessage) → string`).
- Cover:
  - round-trip `Decrypt(Encrypt(text)) == text` across a matrix of inputs (plain ASCII,
    Unicode, multi-line, empty, near the 100000-char limit);
  - tampered ciphertext (flipped/truncated `Data`) → error, no plaintext;
  - corrupted key material (truncated `EncryptedKey`, garbage `PrivateKey`) → error;
  - empty and maximum-size (100000-char) inputs round-trip.
- Use the standard library `testing` package only — no new dependency. `secret_assistant`
  is a stdlib-first module and its `go.mod` has no `testify`; adding it just for these
  simple table tests would violate the "agree on new dependencies" rule in `AGENTS.md`.
- No production-code change, no gRPC/HTTP contract change, no CI change (`go test ./...`
  already runs in CI).

## Capabilities

### New Capabilities

- `crypto-tests`: the crypto module's correctness properties (round-trip,
  tamper-rejection, boundary inputs) are pinned by unit tests.

## Impact

- `secret_assistant/internal/encrypt/encrypt_test.go`,
  `secret_assistant/internal/decrypt/decrypt_test.go` (new test files only). No
  production code, no `go.mod` change.

## Sequencing

Must land before `secure-key-storage`: it pins the current hybrid RSA+AES behaviour, and
`secure-key-storage` then extends these tests to the new opaque-blob/service-key signature
(including the "wrong service key" case, which only becomes testable once a service key
exists).
