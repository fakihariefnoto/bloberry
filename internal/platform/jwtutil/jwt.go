package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"uid"`
	TenantID string `json:"tid"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Issuer signs HS256 JWTs with the algorithm pinned (domains.md §4.2).
type Issuer struct {
	secret []byte
}

func NewIssuer(secret []byte) *Issuer { return &Issuer{secret: secret} }

func (i *Issuer) Issue(userID string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(i.secret)
}

func (i *Issuer) Verify(token string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256"}))
	claims := &Claims{}
	parsed, err := parser.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return i.secret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("jwt: invalid")
	}
	return claims, nil
}
