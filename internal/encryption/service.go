package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

func (s *service) GenerateUserEncryptionKey() ([]byte, error) {
	key := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("%w: %v", EncryptionFailed, err)
	}

	return key, nil
}

func (s *service) EncryptUserKey(userKey []byte) ([]byte, error) {
	return s.Encrypt(userKey, s.masterKey)
}

func (s *service) DecryptUserKey(encrypted []byte) ([]byte, error) {
	return s.Decrypt(encrypted, s.masterKey)
}

func (s *service) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", EncryptionFailed, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", EncryptionFailed, err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("%w: %v", EncryptionFailed, err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (s *service) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", DecryptionFailed, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", DecryptionFailed, err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, CiphertextTooShort
	}

	nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", DecryptionFailed, err)
	}

	return plaintext, nil
}

func (s *service) Hash(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}
