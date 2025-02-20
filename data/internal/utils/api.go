package utils

import "golang.org/x/crypto/bcrypt"

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
