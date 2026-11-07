// Package keys loads the service-held RSA private key used for envelope
// encryption. The private key never leaves this process: only the public key
// derived from it is handed to the encrypt path.
package keys

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// Load reads the RSA-2048 private key from the PEM file at path. It accepts
// both PKCS#1 ("RSA PRIVATE KEY") and PKCS#8 ("PRIVATE KEY") encodings and
// rejects anything that is not a valid RSA-2048 private key.
func Load(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("no PEM block found in key file")
	}

	key, err := parseBlock(block)
	if err != nil {
		return nil, err
	}

	if err := key.Validate(); err != nil {
		return nil, fmt.Errorf("invalid RSA private key: %w", err)
	}
	if key.N.BitLen() != 2048 {
		return nil, fmt.Errorf("RSA private key must be 2048 bits, got %d", key.N.BitLen())
	}

	return key, nil
}

func parseBlock(block *pem.Block) (*rsa.PrivateKey, error) {
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse PKCS#1 private key: %w", err)
		}
		return key, nil
	case "PRIVATE KEY":
		parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse PKCS#8 private key: %w", err)
		}
		key, ok := parsed.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key is not an RSA private key (got %T)", parsed)
		}
		return key, nil
	default:
		return nil, fmt.Errorf("unsupported PEM block type %q", block.Type)
	}
}
