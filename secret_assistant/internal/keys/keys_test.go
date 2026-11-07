package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
)

// writeKeyFile writes data to a fresh temp file and returns its path.
func writeKeyFile(t *testing.T, name string, data []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	return path
}

// pkcs1PEM marshals key as a PKCS#1 "RSA PRIVATE KEY" PEM block.
func pkcs1PEM(t *testing.T, key *rsa.PrivateKey) []byte {
	t.Helper()
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
}

// pkcs8PEM marshals key as a PKCS#8 "PRIVATE KEY" PEM block.
func pkcs8PEM(t *testing.T, key *rsa.PrivateKey) []byte {
	t.Helper()

	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("x509.MarshalPKCS8PrivateKey() error = %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
}

func TestLoadPKCS1(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey() error = %v", err)
	}

	loaded, err := Load(writeKeyFile(t, "rsa.pem", pkcs1PEM(t, key)))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.N.Cmp(key.N) != 0 {
		t.Error("Load() returned a different key")
	}
}

func TestLoadPKCS8(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey() error = %v", err)
	}

	loaded, err := Load(writeKeyFile(t, "rsa.pem", pkcs8PEM(t, key)))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.N.Cmp(key.N) != 0 {
		t.Error("Load() returned a different key")
	}
}

func TestLoadMissingFile(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "missing.pem")); err == nil {
		t.Error("Load() with missing file: expected error, got nil")
	}
}

func TestLoadNonPEMGarbage(t *testing.T) {
	if _, err := Load(writeKeyFile(t, "garbage.pem", []byte("this is not a PEM file"))); err == nil {
		t.Error("Load() with non-PEM garbage: expected error, got nil")
	}
}

func TestLoadRejectsNon2048BitKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("rsa.GenerateKey() error = %v", err)
	}

	if _, err := Load(writeKeyFile(t, "small.pem", pkcs1PEM(t, key))); err == nil {
		t.Error("Load() with 1024-bit key: expected error, got nil")
	}
}
