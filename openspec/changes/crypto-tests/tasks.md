# Tasks: Crypto module tests

Role ownership: `backend` = `secret_assistant/**` (unit tests live next to the code and
belong to the code's owner).

## 1. Unit tests for the crypto module (backend)

- [ ] 1.1 Add `secret_assistant/internal/encrypt/encrypt_test.go` covering `Encrypt`
  output shape (non-empty `EncryptedKey` and `Data`, fresh `PrivateKey` per message) and
  the round-trip matrix of inputs (ASCII, Unicode, multi-line, empty, 100000 chars).
- [ ] 1.2 Add `secret_assistant/internal/decrypt/decrypt_test.go` covering tampered/truncated
  `Data`, truncated `EncryptedKey`, and garbage `PrivateKey` all returning an error.
- [ ] 1.3 Use only the standard library `testing` package (no `testify`, no new
  dependency).
- [ ] 1.4 Verify `cd secret_assistant && go test ./... && go vet ./...` green and that
  the packages are no longer reported as "[no test files]".
