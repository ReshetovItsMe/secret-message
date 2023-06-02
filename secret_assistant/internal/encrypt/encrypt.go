package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
)

type EncryptedMessage struct {
	EncryptedKey []byte
	Data         []byte
}

// Generates AES key with random length (16, 24 or 32 bytes)
func generateAESKey(size int) []byte {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return key
}

// Encrypt Text using AES
func encryptText(text string, aesKey []byte) ([]byte, error) {
	key := generateAESKey(32)

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

// AES key encryption using RSA
func encryptAESKey(aesKey []byte) ([]byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	publicKey := &privateKey.PublicKey
	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, aesKey, nil)
	if err != nil {
		return nil, err
	}
	return encryptedAESKey, err
}

func Encrypt(text string) (*EncryptedMessage, error) {
	aesKey := generateAESKey(32)

	encryptedData, err := encryptText(text, aesKey)
	if err != nil {
		return nil, err
	}

	encryptedAESKey, err := encryptAESKey(aesKey)
	if err != nil {
		return nil, err
	}

	return &EncryptedMessage{encryptedAESKey, encryptedData}, nil

}
