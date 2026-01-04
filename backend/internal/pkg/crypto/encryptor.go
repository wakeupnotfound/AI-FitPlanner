package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type Encryptor struct {
	key []byte
}

func NewEncryptor(secretKey string) (*Encryptor, error) {
	// 确保密钥长度为32字节(256位)
	key := secretKey
	if len(key) < 32 {
		// 如果密钥太短,用0填充
		key = fmt.Sprintf("%-32s", key)
	} else if len(key) > 32 {
		// 如果密钥太长,截取前32字节
		key = key[:32]
	}

	return &Encryptor{
		key: []byte(key),
	}, nil
}

func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("创建cipher失败: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM失败: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成nonce失败: %w", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// 返回十六进制字符串
	return hex.EncodeToString(ciphertext), nil
}

func (e *Encryptor) Decrypt(ciphertextHex string) (string, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", fmt.Errorf("解码十六进制字符串失败: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("创建cipher失败: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM失败: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("密文长度不足")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("解密失败: %w", err)
	}

	return string(plaintext), nil
}

// 快速加密/解密函数
func Encrypt(plaintext, secretKey string) (string, error) {
	encryptor, err := NewEncryptor(secretKey)
	if err != nil {
		return "", err
	}
	return encryptor.Encrypt(plaintext)
}

func Decrypt(ciphertext, secretKey string) (string, error) {
	encryptor, err := NewEncryptor(secretKey)
	if err != nil {
		return "", err
	}
	return encryptor.Decrypt(ciphertext)
}
