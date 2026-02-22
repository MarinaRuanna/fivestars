package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const defaultExpiration = 24 * 7 * time.Hour // 7 dias

// Claims contém os claims do JWT (Subject = user_id).
type Claims struct {
	jwt.RegisteredClaims
}

// NewToken gera um JWT para o userID com o secret dado.
func NewToken(userID, secret string, exp time.Duration) (string, error) {
	if exp == 0 {
		exp = defaultExpiration
	}
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("jwt sign: %w", err)
	}
	return signed, nil
}

// ParseToken valida o token e retorna o userID (sub). Secret deve ser o mesmo usado na assinatura.
func ParseToken(tokenString, secret string) (userID string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	return claims.Subject, nil
}
