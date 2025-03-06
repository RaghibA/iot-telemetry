package utils

import (
    "crypto/rand"
    "encoding/hex"
)

// GenerateAPIKey generates an API key for use by a device.
// Params: None
// Returns:
// - string: the generated API key string
// - error: error if any occurred during the key generation
func GenerateAPIKey() (string, error) {
    bytes := make([]byte, 32) // 256-bit key
    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}
