package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// JWTAdapter implements TokenGenerator and TokenValidator using JWT library
type JWTAdapter struct {
	secret     string
	expiration int // in hours
}

// NewJWTAdapter creates a new JWT adapter
func NewJWTAdapter(secret string, expiration int) *JWTAdapter {
	return &JWTAdapter{
		secret:     secret,
		expiration: expiration,
	}
}

// Generate implements TokenGenerator interface
func (a *JWTAdapter) Generate(userID int64, email string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(a.expiration) * time.Hour)
	claims := &jwtClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validate implements TokenValidator interface
func (a *JWTAdapter) Validate(tokenString string) (*domainuser.TokenClaims, error) {
	claims := &jwtClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return &domainuser.TokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil
}

// jwtClaims represents JWT claims (internal implementation detail)
type jwtClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

