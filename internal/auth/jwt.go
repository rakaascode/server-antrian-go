package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Role     string `json:"role"`
	CabangID *uint  `json:"cabang_id,omitempty"` // hanya untuk admin
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, email, username, role string, cabangID *uint) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Role:     role,
		CabangID: cabangID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
