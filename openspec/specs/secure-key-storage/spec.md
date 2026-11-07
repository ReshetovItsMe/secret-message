# secure-key-storage Specification

## Purpose
Guarantees that a one-time message's stored payload contains no private key material, so
that read access to the storage layer (Redis) alone cannot reveal any plaintext.

## Requirements

### Requirement: No private key material in stored payloads
The payload stored for a one-time message SHALL contain ciphertext and a wrapped data
key only, and SHALL NOT contain any private key material, so that read access to the
storage system (Redis) alone does not reveal the plaintext.

#### Scenario: Redis compromise does not reveal plaintext
- **WHEN** an attacker gains full read access to Redis (all `msg:<uuid>` values)
- **THEN** no stored value contains a usable decryption key, and the attacker cannot
  recover any plaintext without also obtaining the crypto service's private key

#### Scenario: Payload is opaque to the gateway
- **WHEN** the gateway stores or loads a message
- **THEN** it handles the payload as an opaque byte blob and does not parse, name, or
  expose any key or ciphertext field

### Requirement: Service-held envelope key
The crypto service SHALL hold a single long-lived key pair whose private part is
supplied only to the crypto service (via environment-selected PEM file or mounted
secret) and SHALL use it to wrap (encrypt) and unwrap (decrypt) the per-message data
key.

#### Scenario: Private key is not exposed through the API
- **WHEN** the gateway calls the crypto service to encrypt or decrypt
- **THEN** the private key is never returned in any response, error, or log

#### Scenario: Encrypt then decrypt round-trips
- **WHEN** a plaintext is encrypted to a blob and the same blob is later decrypted with
  the same service key
- **THEN** the original plaintext is recovered exactly

#### Scenario: Wrong service key fails to decrypt
- **WHEN** a blob produced under one service key is decrypted with a different service
  key
- **THEN** decryption fails with an error and no plaintext is returned

### Requirement: Versioned opaque blob format
The encrypted payload SHALL be a single self-describing, versioned binary blob that the
crypto service can unambiguously parse without out-of-band metadata.

#### Scenario: Tampered blob is rejected
- **WHEN** any bytes of a stored blob are altered before decryption
- **THEN** decryption fails with an error (authentication failure) rather than returning
  corrupted plaintext
