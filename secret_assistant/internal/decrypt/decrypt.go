package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
)

func Decrypt(text []byte) (string, error) {
	key := "my32digitkey12345678901234567890"

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(text) < nonceSize {
		return "", err
	}

	nonce, ciphertext := text[:nonceSize], text[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext[:]), nil
}
