package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"io"

	m "github.com/ReshetovItsMe/one-time-messaging-exchange-be/pkg/messages"
)

// Generates AES key with random length (16, 24 or 32 bytes)
func generateAESKey(size int) ([]byte, error) {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Encrypt Text using AES
func encryptText(text string, aesKey []byte) ([]byte, error) {
	block, aesErr := aes.NewCipher([]byte(aesKey))
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

// AES key encryption using RSA
func encryptAESKey(aesKey []byte) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey
	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, aesKey, nil)
	if err != nil {
		return nil, nil, err
	}
	privateKeyBytes, _ := json.Marshal(privateKey)

	return encryptedAESKey, privateKeyBytes, err
}

func Encrypt(text string) (*m.EncryptedMessage, error) {
	aesKey, err := generateAESKey(32)
	if err != nil {
		return nil, err
	}

	encryptedData, err := encryptText(text, aesKey)
	if err != nil {
		return nil, err
	}

	encryptedAESKey, privateKey, err := encryptAESKey(aesKey)
	if err != nil {
		return nil, err
	}

	output := &m.EncryptedMessage{PrivateKey: privateKey, EncryptedKey: encryptedAESKey, Data: encryptedData}

	return output, nil
}
