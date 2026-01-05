package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ai-fitness-planner/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims with user information
type Claims struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
	Type      string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTManager interface defines methods for JWT token management
type JWTManager interface {
	GenerateAccessToken(userID int64, username string) (string, error)
	GenerateRefreshToken(userID int64, username string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshAccessToken(refreshToken string) (string, error)
}

// DefaultJWTManager implements the JWTManager interface
type DefaultJWTManager struct {
	secret             string
	accessTokenExpire  time.Duration
	refreshTokenExpire time.Duration
}

// NewJWTManager creates a new JWT manager with configuration
func NewJWTManager(secret string, accessExpire, refreshExpire time.Duration) JWTManager {
	return &DefaultJWTManager{
		secret:             secret,
		accessTokenExpire:  accessExpire,
		refreshTokenExpire: refreshExpire,
	}
}

// NewJWTManagerFromConfig creates a new JWT manager from global config
func NewJWTManagerFromConfig() JWTManager {
	jwtConfig := config.GlobalConfig.JWT
	return NewJWTManager(
		jwtConfig.Secret,
		jwtConfig.AccessTokenExpire,
		jwtConfig.RefreshTokenExpire,
	)
}

// GenerateAccessToken generates a new access token for the user
func (m *DefaultJWTManager) GenerateAccessToken(userID int64, username string) (string, error) {
	sessionID := generateSessionID()

	claims := Claims{
		UserID:    userID,
		Username:  username,
		SessionID: sessionID,
		Type:      "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ai-fitness-planner",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a new refresh token for the user
func (m *DefaultJWTManager) GenerateRefreshToken(userID int64, username string) (string, error) {
	sessionID := generateSessionID()

	claims := Claims{
		UserID:    userID,
		Username:  username,
		SessionID: sessionID,
		Type:      "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ai-fitness-planner",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a token and returns its claims
func (m *DefaultJWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RefreshAccessToken generates a new access token from a valid refresh token
func (m *DefaultJWTManager) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Verify it's a refresh token
	if claims.Type != "refresh" {
		return "", fmt.Errorf("token is not a refresh token")
	}

	// Generate new access token with the same session ID
	accessClaims := Claims{
		UserID:    claims.UserID,
		Username:  claims.Username,
		SessionID: claims.SessionID,
		Type:      "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ai-fitness-planner",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	newAccessToken, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return newAccessToken, nil
}

// generateSessionID generates a unique session ID using crypto/rand
func generateSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("session_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
