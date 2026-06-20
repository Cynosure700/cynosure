package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// New 返回带指定前缀的随机抗碰撞 ID，
// 格式为 "<prefix>_<24 位十六进制字符>"。
func New(prefix string) string {
	return prefix + "_" + Hex()
}

// Hex 返回一个随机的 24 字符十六进制字符串（由 12 个随机字节生成）。
func Hex() string {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "000000000000000000000000"
	}
	return hex.EncodeToString(buf)
}

// UUID 返回一个随机的 RFC 4122 版本 4 UUID 字符串。
func UUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "00000000-0000-4000-8000-000000000000"
	}
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16])
}
