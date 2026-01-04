package jwt

import (
	"fmt"
	"github.com/ai-fitness-planner/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64, username string) (accessToken, refreshToken string, err error) {
	jwtConfig := config.GlobalConfig.JWT

	// Access Token Claims
	accessClaims := Claims{
		UserID:    userID,
		Username:  username,
		SessionID: generateSessionID(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.AccessTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ai-fitness-planner",
		},
	}

	// Refresh Token Claims
	refreshClaims := Claims{
		UserID:    userID,
		Username:  username,
		SessionID: accessClaims.SessionID, // 使用相同的sessionID
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.RefreshTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ai-fitness-planner",
		},
	}

	// 生成Access Token
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", "", fmt.Errorf("生成access token失败: %w", err)
	}

	// 生成Refresh Token
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", "", fmt.Errorf("生成refresh token失败: %w", err)
	}

	return accessToken, refreshToken, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	jwtConfig := config.GlobalConfig.JWT

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析token失败: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的token")
}

func RefreshAccessToken(refreshToken string) (newAccessToken string, err error) {
	claims, err := ParseToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("刷新token无效: %w", err)
	}

	jwtConfig := config.GlobalConfig.JWT

	// 生成新的Access Token
	accessClaims := Claims{
		UserID:    claims.UserID,
		Username:  claims.Username,
		SessionID: claims.SessionID, // 保持相同的sessionID
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.AccessTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ai-fitness-planner",
		},
	}

	newAccessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", fmt.Errorf("生成新的access token失败: %w", err)
	}

	return newAccessToken, nil
}

func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}
