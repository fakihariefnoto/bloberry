package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// Envelope encrypts/decrypts storage-backend credentials with a key from the
// environment (TRD R7). The key never lives in MongoDB.
type Envelope struct {
	key []byte
}

func NewEnvelope(key string) (*Envelope, error) {
	if len(key) < 16 {
		return nil, errors.New("credential encryption key must be at least 16 bytes")
	}
	// SHA-256 the key so any length works and the key is fixed-size for AES-256.
	sum := sha256Sum([]byte(key))
	return &Envelope{key: sum[:]}, nil
}

func sha256Sum(b []byte) [32]byte {
	return sha256.Sum256(b)
}

func (e *Envelope) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (e *Envelope) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ct := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ct, nil)
}

func (e *Envelope) EncryptString(plaintext string) (string, error) {
	b, err := e.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// NewEnvelopeOrPanic is a wiring convenience: the encryption key is required
// at boot (TRD R7), so a bad key is a startup error, not a runtime one.
func NewEnvelopeOrPanic(key string) *Envelope {
	e, err := NewEnvelope(key)
	if err != nil {
		panic(err)
	}
	return e
}

func (e *Envelope) DecryptString(encoded string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	out, err := e.Decrypt(b)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
