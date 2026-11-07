package encrypt

import (
	"strings"
	"testing"

	cryptography "github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/decrypt"
)

func TestEncryptOutputShape(t *testing.T) {
	msg, err := Encrypt("hello")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if len(msg.EncryptedKey) == 0 {
		t.Error("Encrypt() EncryptedKey is empty")
	}
	if len(msg.Data) == 0 {
		t.Error("Encrypt() Data is empty")
	}
	if len(msg.PrivateKey) == 0 {
		t.Error("Encrypt() PrivateKey is empty")
	}
}

func TestEncryptFreshPrivateKeyPerMessage(t *testing.T) {
	first, err := Encrypt("one")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	second, err := Encrypt("two")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if string(first.PrivateKey) == string(second.PrivateKey) {
		t.Error("Encrypt() should generate a fresh PrivateKey for every message")
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "ascii", text: "the quick brown fox"},
		{name: "unicode", text: "Привет, мир! 你好，世界 🌍"},
		{name: "multi-line", text: "line one\nline two\nline three\n"},
		{name: "empty", text: ""},
		{name: "max-100000", text: strings.Repeat("x", 100000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := Encrypt(tt.text)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			got, err := cryptography.Decrypt(msg)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			if got != tt.text {
				t.Fatalf("round-trip mismatch: got %q (%d bytes), want %d bytes", got, len(got), len(tt.text))
			}
		})
	}
}
