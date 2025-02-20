package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

/**
 * Generates api key for use by device
 *
 * @output (string, error): api key string or err
 */
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32) // 256bit key
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

/**
 * Hashes api key for storage in db
 *
 * @params key string: api key string
 * @output (string, error): hashed key or err
 */
func HashAPIKey(key string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

/**
 * Verifies hashed api key against key
 *
 * @params (key string, hashedKey string): api key string, hashed apiKey string from db
 * @output bool: true if keys match
 */
func VerifyAPIKey(key string, hashedKey string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedKey), []byte(key))
	return err == nil
}
