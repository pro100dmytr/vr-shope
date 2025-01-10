package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

func HashPassword(password string) (string, string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", "", err
	}

	hash := sha256.New()
	hash.Write(salt)
	hash.Write([]byte(password))
	hashedPassword := hash.Sum(nil)

	return hex.EncodeToString(hashedPassword), hex.EncodeToString(salt), nil
}

func CheckPassword(password, storedHash, storedSalt string) (bool, error) {
	storedHashBytes, err := hex.DecodeString(storedHash)
	if err != nil {
		return false, fmt.Errorf("failed to decode stored hash: %w", err)
	}

	storedSaltBytes, err := hex.DecodeString(storedSalt)
	if err != nil {
		return false, fmt.Errorf("failed to decode stored salt: %w", err)
	}

	hash := sha256.New()
	hash.Write(storedSaltBytes)
	hash.Write([]byte(password))
	computedHash := hash.Sum(nil)

	return bytes.Equal(computedHash, storedHashBytes), nil
}
