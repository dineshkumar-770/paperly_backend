package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomString() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}