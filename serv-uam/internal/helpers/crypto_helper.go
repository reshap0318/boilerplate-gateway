package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomString generates a cryptographically secure random string.
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashString hashes a string using bcrypt.
func HashString(str string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyString verifies a string against its hash.
func VerifyString(str, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}

// GenerateNumericOTP generates a random zero-padded numeric code of the given length.
func GenerateNumericOTP(length int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", length, n), nil
}

// HashToken deterministically hashes a high-entropy token (e.g. a password-reset token) for
// storage/lookup. Unlike HashString (bcrypt: salted, slow, for low-entropy secrets like
// passwords), this must produce the same output for the same input to be usable in a
// WHERE token = ? lookup — bcrypt can't do that since its salt makes every hash of the same
// input different.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
