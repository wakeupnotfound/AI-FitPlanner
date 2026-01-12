package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
)

// Encryptor interface defines methods for encryption and decryption
type Encryptor interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// AESEncryptor implements Encryptor interface using AES-256-GCM
type AESEncryptor struct {
	key []byte
}

// NewEncryptor creates a new AES-256-GCM encryptor with the given secret key
// Uses PBKDF2 to derive a proper 32-byte key from the input
func NewEncryptor(secretKey string) (*AESEncryptor, error) {
	// Reject empty or very short keys
	if len(secretKey) < 8 {
		return nil, fmt.Errorf("secret key must be at least 8 characters")
	}

	// Use a deterministic salt derived from the secret key
	// This ensures the same secret key always produces the same encryption key
	hash := sha256.Sum256([]byte(secretKey + "salt"))
	salt := hash[:16] // Use first 16 bytes as salt

	// Derive 32-byte key using PBKDF2 with 100,000 iterations
	key := pbkdf2.Key([]byte(secretKey), salt, 100000, 32, sha256.New)

	return &AESEncryptor{
		key: key,
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
// Expects IV to be prepended to the ciphertext
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
