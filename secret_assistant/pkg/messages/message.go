package messages

import "crypto/rsa"

type EncryptedMessage struct {
	PrivateKey   *rsa.PrivateKey
	EncryptedKey []byte
	Data         []byte
}
