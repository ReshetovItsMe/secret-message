package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func Encrypt(text string) ([]byte, error) {
	key := "my32digitkey12345678901234567890"

	block, aesErr := aes.NewCipher([]byte(key))
	if aesErr != nil {
		fmt.Println("ERROR HERE")
		fmt.Println(aesErr)
		return nil, aesErr
	}
	gcm, cipherErr := cipher.NewGCM(block)
	if cipherErr != nil {
		fmt.Println("ERROR HERE 2")
		fmt.Println(cipherErr)
		return nil, cipherErr
	}

	nonce := make([]byte, gcm.NonceSize())

	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("ERROR HERE 3")
		fmt.Println(err)
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)

	return ciphertext, nil
}
