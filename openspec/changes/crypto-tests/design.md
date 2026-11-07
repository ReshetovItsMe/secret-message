# Design: Crypto module tests

## Context

`secret_assistant` has zero tests. `internal/encrypt.Encrypt` and
`internal/decrypt.Decrypt` are pure functions over `*messages.EncryptedMessage`, with no
I/O or external services. The module is stdlib-first and its `go.mod` has no `testify`
(unlike the gateway). See proposal.md — Why.

## Goals / Non-Goals

**Goals:**
- Pin the current hybrid RSA-OAEP + AES-GCM behaviour before `secure-key-storage` changes
  the signature.
- Cover the security-relevant failure modes: round-trip, tampered ciphertext, corrupted
  key material, boundary inputs.

**Non-Goals:**
- No production-code or contract change.
- No new dependency.
- No rewrite of the crypto implementation.

## Decisions

- **Standard library `testing` only, no `testify`.** `secret_assistant` has no `testify`
  and `AGENTS.md` requires agreeing on new dependencies; plain `testing` with table-driven
  subtests and `bytes.Equal`/error checks is sufficient for these cases.
- **"Wrong key" maps to key-material corruption today.** The current code generates a
  fresh RSA keypair per `Encrypt`, so there is no service key to mismatch. The "wrong key"
  requirement is therefore exercised as a truncated `EncryptedKey` or garbage `PrivateKey`
  on the existing `EncryptedMessage` shape; the true "wrong service key" case lands with
  `secure-key-storage`.
- **Tampering targets `Data`.** AES-GCM authentication is the property under test: flip
  one byte or truncate the ciphertext and assert `Decrypt` returns an error (never a
  wrong or empty plaintext). Tests copy the `[]byte` fields before mutating so one test's
  mutation cannot leak into another.
- **Boundary values are empty and the 100000-char maximum** (the handler's validation
  ceiling). A single round-trip case per extreme is enough.
- **Keep the round-trip matrix small.** Each `Encrypt` generates an RSA-2048 keypair
  (~hundreds of ms), so the round-trip table uses a handful of inputs; the tamper/corruption
  cases reuse one `Encrypt` result instead of re-encrypting per assertion.

## Risks / Trade-offs

- [Tests are coupled to the current signature] → accepted: `secure-key-storage` (task
  group 4) explicitly extends them to the new signature, and the round-trip property is
  signature-independent.
- [RSA keygen makes tests slow] → mitigated by a small round-trip matrix and reuse of a
  single ciphertext for failure cases.
