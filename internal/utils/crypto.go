package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Hash(text string) string {
	bcryptBytes, _ := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	return string(bcryptBytes)
}

func CompareHashAndText(hash string, text string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))

	return err == nil
}

func CreateToken[T any](claims T) (string, error) {
	// Define a secret key (use a secure key in production)
	secretKey := []byte(configurations.Config.JWT.Secret)

	// construct claims
	claimsObject := jwt.MapClaims{}
	ObjectMapper(claims, &claimsObject)

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsObject)

	// Sign the token with the secret key
	tokenString, _ := token.SignedString(secretKey)

	return tokenString, nil
}

func EncryptAES(plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(configurations.Config.Application.Secret))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := make([]byte, 12) // GCM nonce size is 12 bytes
	io.ReadFull(rand.Reader, iv)

	aead, _ := cipher.NewGCM(block)
	ciphertext = append(iv, aead.Seal(nil, iv, []byte(plaintext), nil)...)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptAES(ciphertext string) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, _ := aes.NewCipher([]byte(configurations.Config.Application.Secret))

	if len(ciphertextBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertextBytes[:12] // GCM nonce size is 12 bytes
	ciphertextBytes = ciphertextBytes[12:]

	aead, _ := cipher.NewGCM(block)
	plaintext, _ := aead.Open(nil, iv, ciphertextBytes, nil)

	return string(plaintext), nil
}
