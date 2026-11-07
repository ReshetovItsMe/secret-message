package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/encrypt"
)

// newTestKey returns a fresh RSA-2048 keypair.
func newTestKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey() error = %v", err)
	}
	return key
}

// encryptWith seals plaintext with key's public key and returns the blob.
func encryptWith(t *testing.T, key *rsa.PrivateKey, plaintext string) []byte {
	t.Helper()

	blob, err := encrypt.Encrypt(&key.PublicKey, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	return blob
}

func TestDecryptNilPrivateKey(t *testing.T) {
	if _, err := Decrypt(nil, []byte{1, 0, 0, 0}); err == nil {
		t.Error("Decrypt() with nil private key: expected error, got nil")
	}
}

func TestDecryptEmptyBlob(t *testing.T) {
	key := newTestKey(t)

	if _, err := Decrypt(key, nil); err == nil {
		t.Error("Decrypt() with nil blob: expected error, got nil")
	}
	if _, err := Decrypt(key, []byte{}); err == nil {
		t.Error("Decrypt() with empty blob: expected error, got nil")
	}
}

func TestDecryptWrongServiceKey(t *testing.T) {
	key := newTestKey(t)
	wrongKey := newTestKey(t)

	blob := encryptWith(t, key, "a secret only the right service key can open")

	if _, err := Decrypt(wrongKey, blob); err == nil {
		t.Error("Decrypt() with wrong service key: expected error, got nil")
	}
}

func TestDecryptTamperedCiphertext(t *testing.T) {
	key := newTestKey(t)

	blob := encryptWith(t, key, "some secret message that is long enough to tamper with")

	// Flip the last byte (part of the AES-GCM ciphertext + tag).
	tampered := append([]byte(nil), blob...)
	tampered[len(tampered)-1] ^= 0xff

	if _, err := Decrypt(key, tampered); err == nil {
		t.Error("Decrypt() with tampered ciphertext: expected error, got nil")
	}
}

func TestDecryptTamperedWrappedKey(t *testing.T) {
	key := newTestKey(t)

	blob := encryptWith(t, key, "some secret message")

	// Flip the first byte of the wrapped AES key (byte 1, right after the
	// version byte).
	tampered := append([]byte(nil), blob...)
	tampered[1] ^= 0xff

	if _, err := Decrypt(key, tampered); err == nil {
		t.Error("Decrypt() with tampered wrapped key: expected error, got nil")
	}
}

func TestDecryptUnknownVersion(t *testing.T) {
	key := newTestKey(t)

	blob := encryptWith(t, key, "some secret message")

	bad := append([]byte(nil), blob...)
	bad[0] = 0x7f

	if _, err := Decrypt(key, bad); err == nil {
		t.Error("Decrypt() with unknown blob version: expected error, got nil")
	}
}

func TestDecryptTruncatedBlob(t *testing.T) {
	key := newTestKey(t)

	blob := encryptWith(t, key, "some secret message that is long enough to truncate")

	for _, n := range []int{0, 1, 5, 200, len(blob) - 1} {
		if _, err := Decrypt(key, blob[:n]); err == nil {
			t.Errorf("Decrypt() with blob truncated to %d bytes: expected error, got nil", n)
		}
	}
}
