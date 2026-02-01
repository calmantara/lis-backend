package utils

import (
	"testing"
	"time"

	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T) {
	text := "password123"
	hashedText := Hash(text)

	assert.NotEmpty(t, hashedText, "Hashed text should not be empty")

	err := bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(text))
	assert.NoError(t, err, "Hashed text should match the original text")
}

func TestCompareHashAndText(t *testing.T) {
	text := "password123"
	hashedText := Hash(text)

	assert.True(t, CompareHashAndText(hashedText, text), "Hash and text should match")
	assert.False(t, CompareHashAndText(hashedText, "wrongpassword"), "Hash and incorrect text should not match")
}

func TestCreateTokenWithUserSessionClaim(t *testing.T) {
	configurations.Load()
	// Mock configurations.JWT_SECRET
	configurations.Config.Lock()
	defer configurations.Config.Unlock()

	configurations.Config.JWT.Secret = "test_secret_key"

	// Mock claims
	claims := map[string]any{
		"sub": "session",
		"uid": "123",
		"iss": "test_issuer",
		"aud": "test_audience",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
		"nbf": time.Now().Unix(),
	}

	token, err := CreateToken(claims)
	assert.NoError(t, err, "Token creation should not return an error")
	assert.NotEmpty(t, token, "Token should not be empty")

	// Parse the token to validate it
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(configurations.Config.JWT.Secret), nil
	})
	assert.NoError(t, err, "Token parsing should not return an error")
	assert.True(t, parsedToken.Valid, "Parsed token should be valid")
}

func TestCreateTokenWithUserTokenClaim(t *testing.T) {
	configurations.Load()
	// Mock configurations.JWT_SECRET
	configurations.Config.Lock()
	defer configurations.Config.Unlock()

	configurations.Config.JWT.Secret = "test_secret_key"

	// Mock claims
	claims := map[string]any{
		"sub": "user",
		"uid": "123",
		"usr": "testuser",
		"iss": "test_issuer",
		"aud": "test_audience",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
		"nbf": time.Now().Unix(),
	}

	token, err := CreateToken(claims)
	assert.NoError(t, err, "Token creation should not return an error")
	assert.NotEmpty(t, token, "Token should not be empty")

	// Parse the token to validate it
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(configurations.Config.JWT.Secret), nil
	})
	assert.NoError(t, err, "Token parsing should not return an error")
	assert.True(t, parsedToken.Valid, "Parsed token should be valid")
}

func TestEncryptAESWithEmptyPlaintext(t *testing.T) {
	configurations.Load()
	// Mock configurations.Application.Secret
	configurations.Config.Lock()
	defer configurations.Config.Unlock()

	configurations.Config.Application.Secret = "JzkMNeV5o2nF425kgduBfFjSEi8ah5CL"

	plaintext := ""
	encryptedText, err := EncryptAES(plaintext)

	assert.NoError(t, err, "Encryption should not return an error for empty plaintext")
	assert.NotEmpty(t, encryptedText, "Encrypted text should not be empty for empty plaintext")
}

func TestDecryptAESWithEmptyCiphertext(t *testing.T) {
	configurations.Load()
	// Mock configurations.Application.Secret
	configurations.Config.Lock()
	defer configurations.Config.Unlock()

	configurations.Config.Application.Secret = "JzkMNeV5o2nF425kgduBfFjSEi8ah5CL"

	ciphertext := ""
	decryptedText, err := DecryptAES(ciphertext)

	assert.Error(t, err, "Decryption should return an error for empty ciphertext")
	assert.Empty(t, decryptedText, "Decrypted text should be empty for empty ciphertext")
}

func TestEncryptAESWithSpecialCharacters(t *testing.T) {
	configurations.Load()
	// Mock configurations.Application.Secret
	configurations.Config.Lock()
	defer configurations.Config.Unlock()

	configurations.Config.Application.Secret = "JzkMNeV5o2nF425kgduBfFjSEi8ah5CL"

	plaintext := "Special characters: !@#$%^&*()_+-=[]{}|;:',.<>?/`~"
	encryptedText, err := EncryptAES(plaintext)

	assert.NoError(t, err, "Encryption should not return an error for special characters")
	assert.NotEmpty(t, encryptedText, "Encrypted text should not be empty for special characters")
	assert.NotEqual(t, plaintext, encryptedText, "Encrypted text should not match the plaintext")
}

func TestDecryptAESWithSpecialCharacters(t *testing.T) {
	configurations.Load()
	// Mock configurations.Application.Secret
	configurations.Config.Lock()
	defer configurations.Config.Unlock()

	configurations.Config.Application.Secret = "JzkMNeV5o2nF425kgduBfFjSEi8ah5CL"

	plaintext := "Special characters: !@#$%^&*()_+-=[]{}|;:',.<>?/`~"
	encryptedText, err := EncryptAES(plaintext)
	assert.NoError(t, err, "Encryption should not return an error for special characters")

	decryptedText, err := DecryptAES(encryptedText)
	assert.NoError(t, err, "Decryption should not return an error for special characters")
	assert.Equal(t, plaintext, decryptedText, "Decrypted text should match the original plaintext with special characters")
}

func TestDecryptAESWithModifiedCiphertext(t *testing.T) {
	configurations.Load()
	// Mock configurations.Application.Secret
	configurations.Config.Lock()
	defer configurations.Config.Unlock()

	configurations.Config.Application.Secret = "JzkMNeV5o2nF425kgduBfFjSEi8ah5CL"

	plaintext := "This is a secret message"
	encryptedText, err := EncryptAES(plaintext)
	assert.NoError(t, err, "Encryption should not return an error")

	// Modify the ciphertext
	modifiedCiphertext := encryptedText[:len(encryptedText)-1] + "A"

	decryptedText, err := DecryptAES(modifiedCiphertext)
	assert.Error(t, err, "Decryption should return an error for modified ciphertext")
	assert.Empty(t, decryptedText, "Decrypted text should be empty for modified ciphertext")
}

func TestEncryptAESWithInvalidKeyLength(t *testing.T) {
	configurations.Load()
	// Mock configurations.Application.Secret with an invalid key length
	configurations.Config.Lock()
	defer configurations.Config.Unlock()

	configurations.Config.Application.Secret = "short_key"

	plaintext := "This is a test message"
	encryptedText, err := EncryptAES(plaintext)

	assert.Error(t, err, "Encryption should return an error for invalid key length")
	assert.Empty(t, encryptedText, "Encrypted text should be empty for invalid key length")
}

func TestDecryptAESWithInvalidKeyLength(t *testing.T) {
	configurations.Load()
	// Mock configurations.Application.Secret with an invalid key length
	configurations.Config.Lock()
	defer configurations.Config.Unlock()

	configurations.Config.Application.Secret = "short_key"

	ciphertext := "InvalidCiphertext"
	decryptedText, err := DecryptAES(ciphertext)

	assert.Error(t, err, "Decryption should return an error for invalid key length")
	assert.Empty(t, decryptedText, "Decrypted text should be empty for invalid key length")
}

func TestDecryptAESWithInvalidBase64(t *testing.T) {
	configurations.Load()
	// Mock configurations.Application.Secret
	configurations.Config.Lock()
	defer configurations.Config.Unlock()

	configurations.Config.Application.Secret = "JzkMNeV5o2nF425kgduBfFjSEi8ah5CL"

	invalidBase64 := "InvalidBase64String!!"
	decryptedText, err := DecryptAES(invalidBase64)

	assert.Error(t, err, "Decryption should return an error for invalid base64 input")
	assert.Empty(t, decryptedText, "Decrypted text should be empty for invalid base64 input")
}

func TestDecryptAESWithShortCiphertext(t *testing.T) {
	configurations.Load()
	// Mock configurations.Application.Secret
	configurations.Config.Lock()
	defer configurations.Config.Unlock()

	configurations.Config.Application.Secret = "JzkMNeV5o2nF425kgduBfFjSEi8ah5CL"

	shortCiphertext := "short"
	decryptedText, err := DecryptAES(shortCiphertext)

	assert.Error(t, err, "Decryption should return an error for ciphertext too short")
	assert.Empty(t, decryptedText, "Decrypted text should be empty for ciphertext too short")
}
