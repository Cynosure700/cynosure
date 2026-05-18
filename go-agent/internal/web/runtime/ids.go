package runtime

import (
	"crypto/rand"
	"fmt"
)

func newMessageID() string  { return newID("msg") }
func newToolCallID() string { return newID("tool") }

func newID(prefix string) string {
	buf := make([]byte, 12)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("%s_%x", prefix, buf)
}
