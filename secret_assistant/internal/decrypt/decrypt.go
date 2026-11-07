// Package cryptography opens the opaque versioned blobs produced by the
// encrypt package, unwrapping the ephemeral AES key with the service RSA
// private key and AES-GCM-decrypting the payload.
package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/cryptoblob"
)

// Decrypt opens a blob produced by encrypt.Encrypt using the service private
// key. It rejects an unknown version, a truncated or malformed blob, a key
// that cannot unwrap the AES key (e.g. the wrong service key), and any
// authentication failure from tampering or truncation.
func Decrypt(privateKey *rsa.PrivateKey, blob []byte) (string, error) {
	if privateKey == nil {
		return "", errors.New("decrypt: private key is nil")
	}
	if len(blob) < 1 {
		return "", errors.New("decrypt: empty blob")
	}
	if blob[0] != cryptoblob.Version {
		return "", fmt.Errorf("decrypt: unsupported blob version %d", blob[0])
	}

	rest := blob[1:]

	keySize := privateKey.Size()
	if len(rest) < keySize {
		return "", errors.New("decrypt: blob truncated: missing wrapped key")
	}
	wrappedKey := rest[:keySize]
	rest = rest[keySize:]

	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, wrappedKey, nil)
	if err != nil {
		return "", fmt.Errorf("unwrap AES key: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("new AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(rest) < nonceSize {
		return "", errors.New("decrypt: blob truncated: missing nonce")
	}
	nonce := rest[:nonceSize]
	ciphertext := rest[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt payload: %w", err)
	}

	return string(plaintext), nil
}
