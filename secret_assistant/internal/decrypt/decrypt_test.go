package cryptography

import (
	"testing"

	"github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/encrypt"
	m "github.com/ReshetovItsMe/one-time-messaging-exchange-be/pkg/messages"
)

// validMessage returns a freshly encrypted message to corrupt in each test.
func validMessage(t *testing.T) *m.EncryptedMessage {
	t.Helper()

	msg, err := encrypt.Encrypt("some secret message that is long enough to tamper with")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	return msg
}

func TestDecryptTamperedData(t *testing.T) {
	msg := validMessage(t)

	tampered := *msg
	tampered.Data = append([]byte(nil), msg.Data...)
	tampered.Data[len(tampered.Data)/2] ^= 0xff

	if _, err := Decrypt(&tampered); err == nil {
		t.Error("Decrypt() with tampered Data: expected error, got nil")
	}
}

func TestDecryptTruncatedData(t *testing.T) {
	msg := validMessage(t)

	// Drop the trailing byte of the ciphertext (including the auth tag), but
	// keep the nonce prefix intact so the GCM layer reports the truncation.
	truncated := *msg
	truncated.Data = msg.Data[:len(msg.Data)-1]

	if _, err := Decrypt(&truncated); err == nil {
		t.Error("Decrypt() with truncated Data: expected error, got nil")
	}
}

func TestDecryptTruncatedEncryptedKey(t *testing.T) {
	msg := validMessage(t)

	truncated := *msg
	truncated.EncryptedKey = msg.EncryptedKey[:len(msg.EncryptedKey)-10]

	if _, err := Decrypt(&truncated); err == nil {
		t.Error("Decrypt() with truncated EncryptedKey: expected error, got nil")
	}
}

func TestDecryptGarbagePrivateKey(t *testing.T) {
	msg := validMessage(t)

	// A structurally valid JSON object that is not an RSA private key. A
	// fully non-JSON body is a separate, known robustness gap in Decrypt and
	// is out of scope for these tests (see crypto-tests change).
	garbage := *msg
	garbage.PrivateKey = []byte("{}")

	if _, err := Decrypt(&garbage); err == nil {
		t.Error("Decrypt() with garbage PrivateKey: expected error, got nil")
	}
}
