package config

import (
	"crypto/sha256"
	"fmt"
	"os"
)

// HashFile returns the SHA256 hex digest of a file's contents.
// Returns empty string and error if the file cannot be read.
func HashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

// HashBytes returns the SHA256 hex digest of a byte slice.
func HashBytes(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}
