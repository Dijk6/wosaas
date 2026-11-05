package trace

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
)

// ID 32 字节的 trace id。
type ID [32]byte

func (id ID) String() string {
	return string(id[:])
}

// Generator 生成 trace id。
type Generator struct{}

// NewTraceGenerator 创建 trace id 生成器。
func NewTraceGenerator() *Generator {
	return &Generator{}
}

// NewTraceID 生成新的 trace id。
func (g *Generator) NewTraceID() ID {
	var id ID

	// 生成 32 位十六进制字符串。
	u := uuid.New()
	clean := strings.ReplaceAll(u.String(), "-", "")
	copy(id[:], clean)
	return id
}

// NewRandomID 生成随机十六进制 ID，可用于订单号、请求 ID 等。
func NewRandomID(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return hex.EncodeToString(buf)
}
