package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

// UUID returns a random RFC 4122 version 4 UUID string.
func UUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "00000000-0000-4000-8000-000000000000"
	}
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16])
}
