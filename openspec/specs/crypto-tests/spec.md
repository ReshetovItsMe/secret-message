# crypto-tests Specification

## Purpose
Defines the correctness properties of the crypto module's hybrid RSA-OAEP + AES-GCM
encrypt/decrypt and requires those properties to be pinned by unit tests.

## Requirements

### Requirement: Encrypt then decrypt round-trips
The crypto module SHALL round-trip: for any plaintext the product accepts (including
empty and maximum-size inputs), `Decrypt(Encrypt(text))` SHALL return exactly `text`.

#### Scenario: Round-trip across input shapes
- **WHEN** a plaintext (ASCII, Unicode, multi-line, empty, or 100000 characters) is
  encrypted and the resulting payload is decrypted with its own key material
- **THEN** the original plaintext is recovered exactly

### Requirement: Tampered ciphertext is rejected
The crypto module SHALL reject a payload whose ciphertext has been altered or truncated
with an error, rather than returning corrupted or empty plaintext.

#### Scenario: Altered or truncated ciphertext fails
- **WHEN** any byte of the stored ciphertext (`Data`) is altered, or the ciphertext is
  truncated before decryption
- **THEN** decryption returns an error and no plaintext is produced

### Requirement: Corrupted key material is rejected
The crypto module SHALL reject a payload whose wrapped key or private key material is
corrupted or truncated, returning an error rather than panicking or returning plaintext.

#### Scenario: Corrupted encrypted key or private key fails
- **WHEN** the wrapped AES key (`EncryptedKey`) is truncated, or the private key bytes are
  replaced with garbage
- **THEN** decryption returns an error and no plaintext is produced
