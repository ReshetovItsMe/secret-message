package messages

type EncryptedMessage struct {
	PrivateKey   []byte `json:"privateKey"`
	EncryptedKey []byte `json:"encryptedKey"`
	Data         []byte `json:"data"`
}
