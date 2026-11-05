// Package sessionlib 提供轻量登录态签发与校验。
//
// 当前使用 HMAC-SHA256 签名的自包含 Token，后续可以替换为 JWT/Redis Session。
package sessionlib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// Claims 登录态信息。
type Claims struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Mobile   string `json:"mobile"`
	Role     string `json:"role"`
	SpID     string `json:"sp_id"`
	CpID     string `json:"cp_id"`
	FleetID  string `json:"fleet_id"`
	CityID   int    `json:"city_id"`
	ExpireAt int64  `json:"expire_at"`
}

// Manager Token 管理器。
type Manager struct {
	secret []byte
	ttl    time.Duration
}

// NewManager 创建管理器。
func NewManager(secret string, ttl time.Duration) *Manager {
	if secret == "" {
		secret = "wosaas-dev-secret"
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &Manager{secret: []byte(secret), ttl: ttl}
}

// Sign 签发 token。
func (m *Manager) Sign(claims Claims) (string, error) {
	if claims.ExpireAt == 0 {
		claims.ExpireAt = time.Now().Add(m.ttl).Unix()
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	body := base64.RawURLEncoding.EncodeToString(payload)
	sig := m.signature(body)
	return body + "." + sig, nil
}

// Parse 校验并解析 token。
func (m *Manager) Parse(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, ErrInvalidToken
	}
	if !hmac.Equal([]byte(m.signature(parts[0])), []byte(parts[1])) {
		return nil, ErrInvalidToken
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims := new(Claims)
	if err := json.Unmarshal(payload, claims); err != nil {
		return nil, ErrInvalidToken
	}
	if claims.ExpireAt > 0 && time.Now().Unix() > claims.ExpireAt {
		return nil, ErrExpiredToken
	}
	return claims, nil
}

func (m *Manager) signature(body string) string {
	h := hmac.New(sha256.New, m.secret)
	h.Write([]byte(body))
	return hex.EncodeToString(h.Sum(nil))
}

// Default 全局管理器，由 bootstrap 初始化。
var Default = NewManager("wosaas-dev-secret", 24*time.Hour)

// Init 初始化全局管理器。
func Init(secret string, ttl time.Duration) {
	Default = NewManager(secret, ttl)
}
