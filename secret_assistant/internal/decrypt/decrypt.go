package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	rsa "crypto/rsa"
	"crypto/sha256"

	m "github.com/ReshetovItsMe/one-time-messaging-exchange-be/pkg/messages"
)

// AES key decryption using RSA
func decryptAESKey(privateKey *rsa.PrivateKey, aesKey []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, aesKey, nil)
}

func Decrypt(message *m.EncryptedMessage) (string, error) {
	key, err := decryptAESKey(message.PrivateKey, message.EncryptedKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(message.Data) < nonceSize {
		return "", err
	}

	nonce, ciphertext := message.Data[:nonceSize], message.Data[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext[:]), nil
}
