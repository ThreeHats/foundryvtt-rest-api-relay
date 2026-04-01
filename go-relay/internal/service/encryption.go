package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
)

// EncryptedData holds the result of AES-256-GCM encryption.
type EncryptedData struct {
	Ciphertext string `json:"ciphertext"`
	IV         string `json:"iv"`
	AuthTag    string `json:"authTag"`
}

var hexPattern = regexp.MustCompile(`^[0-9a-fA-F]+$`)

// getEncryptionKey parses the encryption key from hex or base64 format.
func getEncryptionKey(keyEnv string) ([]byte, error) {
	if keyEnv == "" {
		return nil, fmt.Errorf("CREDENTIALS_ENCRYPTION_KEY is not set")
	}

	// Hex encoded (64 chars = 32 bytes)
	if len(keyEnv) == 64 && hexPattern.MatchString(keyEnv) {
		return hex.DecodeString(keyEnv)
	}

	// Base64 encoded (44 chars = 32 bytes)
	decoded, err := base64.StdEncoding.DecodeString(keyEnv)
	if err == nil && len(decoded) == 32 {
		return decoded, nil
	}

	return nil, fmt.Errorf("CREDENTIALS_ENCRYPTION_KEY must be a 32-byte key encoded as hex (64 chars) or base64 (44 chars)")
}

// Encrypt encrypts plaintext using AES-256-GCM.
func Encrypt(plaintext, keyEnv string) (*EncryptedData, error) {
	key, err := getEncryptionKey(keyEnv)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	iv := make([]byte, 12)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("generate IV: %w", err)
	}

	// GCM Seal appends the auth tag to the ciphertext
	sealed := gcm.Seal(nil, iv, []byte(plaintext), nil)

	// Split ciphertext and auth tag (last 16 bytes)
	ciphertext := sealed[:len(sealed)-16]
	authTag := sealed[len(sealed)-16:]

	return &EncryptedData{
		Ciphertext: hex.EncodeToString(ciphertext),
		IV:         hex.EncodeToString(iv),
		AuthTag:    hex.EncodeToString(authTag),
	}, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM.
func Decrypt(ciphertextHex, ivHex, authTagHex, keyEnv string) (string, error) {
	key, err := getEncryptionKey(keyEnv)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}

	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", fmt.Errorf("decode IV: %w", err)
	}

	authTag, err := hex.DecodeString(authTagHex)
	if err != nil {
		return "", fmt.Errorf("decode auth tag: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	// Reconstruct sealed data (ciphertext + auth tag)
	sealed := append(ciphertext, authTag...)

	plaintext, err := gcm.Open(nil, iv, sealed, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	return string(plaintext), nil
}

// IsEncryptionAvailable checks if the encryption key is configured.
func IsEncryptionAvailable(keyEnv string) bool {
	_, err := getEncryptionKey(keyEnv)
	return err == nil
}
