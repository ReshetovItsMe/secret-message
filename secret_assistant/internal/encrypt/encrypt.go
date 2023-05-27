package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func Encrypt(text string) ([]byte, error) {
	key := "my32digitkey12345678901234567890"

	block, aesErr := aes.NewCipher([]byte(key))
	if aesErr != nil {
		return nil, aesErr
	}
	gcm, cipherErr := cipher.NewGCM(block)
	if cipherErr != nil {
		return nil, cipherErr
	}

	nonce := make([]byte, gcm.NonceSize())

	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)

	return ciphertext, nil
}
