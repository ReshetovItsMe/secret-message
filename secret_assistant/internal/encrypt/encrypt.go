// Package encrypt seals plaintext using envelope encryption: the plaintext is
// encrypted with an ephemeral AES-256-GCM key, which is then wrapped with the
// service RSA public key (OAEP/SHA-256). The result is an opaque, versioned
// binary blob that carries no private key material and no JSON.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/cryptoblob"
)

const (
	// aesKeySize is the size in bytes of the ephemeral AES-256 key.
	aesKeySize = 32
)

// Encrypt seals plaintext and returns an opaque, versioned binary blob:
//
//	[version 1B][wrapped AES key 256B][nonce 12B][ciphertext + tag ...]
//
// The AES key is wrapped with the service public key via RSA-OAEP(SHA-256);
// only the matching service private key can unwrap it. No private key material
// is ever produced.
func Encrypt(publicKey *rsa.PublicKey, plaintext string) ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("encrypt: public key is nil")
	}

	aesKey := make([]byte, aesKeySize)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, fmt.Errorf("generate AES key: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("new AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	wrappedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, aesKey, nil)
	if err != nil {
		return nil, fmt.Errorf("wrap AES key: %w", err)
	}

	blob := make([]byte, 0, 1+len(wrappedKey)+len(nonce)+len(ciphertext))
	blob = append(blob, cryptoblob.Version)
	blob = append(blob, wrappedKey...)
	blob = append(blob, nonce...)
	blob = append(blob, ciphertext...)

	return blob, nil
}
