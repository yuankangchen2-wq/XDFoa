package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 载荷
type Claims struct {
	UserID    int64    `json:"user_id"`
	Username  string   `json:"username"`
	TenantID  string   `json:"tenant_id"`
	Roles     []string `json:"roles"`
	OrgPath   string   `json:"org_path"`
	Type      string   `json:"type"` // access / refresh
	jwt.RegisteredClaims
}

// JWT 管理器
type JWTManager struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTManager(secret string, accessTTL, refreshTTL int64) *JWTManager {
	return &JWTManager{
		secret:          []byte(secret),
		accessTokenTTL:  time.Duration(accessTTL) * time.Second,
		refreshTokenTTL: time.Duration(refreshTTL) * time.Second,
	}
}

// GenerateAccessToken 生成 access token
func (m *JWTManager) GenerateAccessToken(userID int64, username, tenantID string, roles []string, orgPath string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		TenantID: tenantID,
		Roles:    roles,
		OrgPath:  orgPath,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// GenerateRefreshToken 生成 refresh token
func (m *JWTManager) GenerateRefreshToken(userID int64, username, tenantID string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		TenantID: tenantID,
		Type:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Parse 解析 token
func (m *JWTManager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
