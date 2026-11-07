package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"strings"
	"testing"

	"github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/cryptoblob"
	cryptography "github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/decrypt"
)

// newTestKey returns a fresh RSA-2048 keypair whose public key is used to seal
// and whose private key is used to open in the round-trip tests.
func newTestKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey() error = %v", err)
	}
	return key
}

func TestEncryptNilPublicKey(t *testing.T) {
	if _, err := Encrypt(nil, "hello"); err == nil {
		t.Error("Encrypt() with nil public key: expected error, got nil")
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := newTestKey(t)

	tests := []struct {
		name string
		text string
	}{
		{name: "empty", text: ""},
		{name: "ascii", text: "the quick brown fox"},
		{name: "unicode", text: "Привет, мир! 你好，世界 🌍"},
		{name: "multi-line", text: "line one\nline two\nline three\n"},
		{name: "max-100000", text: strings.Repeat("x", 100000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blob, err := Encrypt(&key.PublicKey, tt.text)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			got, err := cryptography.Decrypt(key, blob)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			if got != tt.text {
				t.Fatalf("round-trip mismatch: got %d bytes, want %d", len(got), len(tt.text))
			}
		})
	}
}

func TestEncryptProducesOpaqueBlob(t *testing.T) {
	key := newTestKey(t)

	blob, err := Encrypt(&key.PublicKey, "hello, world")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if len(blob) == 0 || blob[0] != cryptoblob.Version {
		t.Fatalf("Encrypt() blob is missing the version byte: %v", blob)
	}
	if strings.Contains(string(blob), "hello, world") {
		t.Error("Encrypt() blob leaks the plaintext")
	}
	if strings.Contains(string(blob), "PRIVATE KEY") {
		t.Error("Encrypt() blob contains private key material")
	}
}

func TestEncryptProducesDistinctBlobs(t *testing.T) {
	key := newTestKey(t)

	first, err := Encrypt(&key.PublicKey, "same plaintext")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	second, err := Encrypt(&key.PublicKey, "same plaintext")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if string(first) == string(second) {
		t.Error("Encrypt() should produce a distinct blob per message (ephemeral AES key + nonce)")
	}
}
