<p align="center">
  <img src="logo.png" />
  <h1 align="center">Secret Message</h1>
<p>

Secret Message is a one-time secret-sharing service: a sender encrypts a message, gets a
unique link, and the recipient can open that link exactly once to read the decrypted
content. Reading the link atomically deletes the stored payload.

The service uses envelope encryption. Each message is encrypted with a fresh AES-256-GCM
data key, which is in turn wrapped by a single RSA-2048 key held by the crypto service.
Only an opaque encrypted blob is stored, and the decryption key never leaves the crypto
service or reaches the data store.

## Features

- **One-time links:** each link can be opened exactly once; reading it deletes the stored
  payload.
- **Encryption at rest:** stored payloads are opaque encrypted blobs with no key material;
  a compromise of the data store (Redis) alone does not reveal message contents.

## How It Works

1. The sender enters a message.
2. The gateway forwards the plaintext to the crypto service, which returns an opaque
   encrypted blob.
3. The blob is stored under a unique id and a one-time link is generated.
4. The recipient opens the link; the blob is fetched and atomically deleted, then decrypted.
5. The plaintext is shown once; after that the message no longer exists.

## Security note

Encryption is server-side, not end-to-end: the gateway and the crypto service process the
plaintext in memory. A compromise of the crypto service (or its private key) can therefore
reveal message contents. What the design does guarantee is that read access to the data
store alone cannot recover any plaintext.

## Contributing

Thank you for considering contributing to this project! If you have any suggestions, bug
reports, or improvements, feel free to open an issue or submit a pull request.
