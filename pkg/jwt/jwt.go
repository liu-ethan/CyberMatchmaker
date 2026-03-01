/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package jwt

import (
	config "CyberMatchmaker/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义 JWT 的 Payload 结构
type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// 从 config 中读取密钥（运行时读取，避免 init 时 AppConfig 为空）
func getSecret() ([]byte, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}
	if cfg.Jwt.Secret == "" {
		return nil, errors.New("jwt secret is empty")
	}
	return []byte(cfg.Jwt.Secret), nil
}

func getConfig() (*config.Config, error) {
	if config.AppConfig == nil {
		return nil, errors.New("config not initialized")
	}
	return config.AppConfig, nil
}

// GenerateToken 生成 Token
func GenerateToken(userID int64) (string, error) {
	cfg, err := getConfig()
	if err != nil {
		return "", err
	}
	// 设置 Token 过期时间
	expirationTime := time.Now().Add(time.Duration(cfg.Jwt.Expire) * time.Hour)

	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cyber_matchmaker",
		},
	}

	// 使用 HS256 算法生成 token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, err := getSecret()
	if err != nil {
		return "", err
	}
	// 使用密钥进行签名并获得完整的字符串 token
	return token.SignedString(secret)
}

// ParseToken 解析并校验 Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	secret, err := getSecret()
	if err != nil {
		return nil, err
	}
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 校验签名算法是否为预期的 HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	// 校验并提取自定义的 Claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
