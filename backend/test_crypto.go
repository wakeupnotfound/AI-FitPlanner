package main

import (
	"fmt"
	"log"

	"github.com/ai-fitness-planner/backend/internal/model"
	"github.com/ai-fitness-planner/backend/internal/pkg/crypto"
)

func main() {
	// Test with the specific API key you provided
	secretKey := "your-secret-key-change-in-production-min-32-chars"

	// Create encryptor
	encryptor, err := crypto.NewEncryptor(secretKey)
	if err != nil {
		log.Fatalf("Failed to create encryptor: %v", err)
	}

	// Test encryption
	plaintext := "Hello, World!"
	encrypted, err := encryptor.Encrypt(plaintext)
	if err != nil {
		log.Fatalf("Failed to encrypt: %v", err)
	}

	// Test decryption
	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		log.Fatalf("Failed to decrypt: %v", err)
	}

	fmt.Printf("Original: %s\n", plaintext)
	fmt.Printf("Encrypted: %s\n", encrypted)
	fmt.Printf("Decrypted: %s\n", decrypted)

	// Create sample API config
	apiConfig := &model.AIAPI{
		Provider:        "wenxin",
		Name:            "Test Wenxin Config",
		APIEndpoint:     "https://aip.baidubce.com",
		APIKeyEncrypted: encrypted,
		Status:          1,
		IsDefault:       true,
	}

	// Test decryption of API key
	decryptedKey, err := encryptor.Decrypt(apiConfig.APIKeyEncrypted)
	if err != nil {
		log.Fatalf("Failed to decrypt API key: %v", err)
	}

	fmt.Printf("Successfully decrypted API key: %s\n", decryptedKey)
	fmt.Printf("Key length: %d\n", len(decryptedKey))
}
