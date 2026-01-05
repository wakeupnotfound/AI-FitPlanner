package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// Encryptor interface defines methods for encryption and decryption
type Encryptor interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// AESEncryptor implements the Encryptor interface using AES-256-GCM
type AESEncryptor struct {
	key []byte
}

// NewEncryptor creates a new AES-256-GCM encryptor with the given secret key
// The key must be exactly 32 bytes for AES-256
func NewEncryptor(secretKey string) (Encryptor, error) {
	// Ensure key length is exactly 32 bytes (256 bits)
	key := secretKey
	if len(key) < 32 {
		// If key is too short, pad with zeros
		key = fmt.Sprintf("%-32s", key)
	} else if len(key) > 32 {
		// If key is too long, truncate to 32 bytes
		key = key[:32]
	}

	return &AESEncryptor{
		key: []byte(key),
	}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM with a random IV
// Returns hex-encoded ciphertext with IV prepended
func (e *AESEncryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce (IV)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data (nonce is prepended to ciphertext)
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return hex-encoded string
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts hex-encoded ciphertext using AES-256-GCM
// Expects the IV to be prepended to the ciphertext
func (e *AESEncryptor) Decrypt(ciphertextHex string) (string, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}

// Encrypt is a convenience function for quick encryption
func Encrypt(plaintext, secretKey string) (string, error) {
	encryptor, err := NewEncryptor(secretKey)
	if err != nil {
		return "", err
	}
	return encryptor.Encrypt(plaintext)
}

// Decrypt is a convenience function for quick decryption
func Decrypt(ciphertext, secretKey string) (string, error) {
	encryptor, err := NewEncryptor(secretKey)
	if err != nil {
		return "", err
	}
	return encryptor.Decrypt(ciphertext)
}
