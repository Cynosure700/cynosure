package idgen

import (
	"crypto/rand"
	"encoding/hex"
)

// New returns a random, collision-resistant ID with the given prefix,
// formatted as "<prefix>_<24 hex chars>".
func New(prefix string) string {
	return prefix + "_" + Hex()
}

// Hex returns a random 24-character hex string (12 random bytes).
func Hex() string {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "000000000000000000000000"
	}
	return hex.EncodeToString(buf)
}
