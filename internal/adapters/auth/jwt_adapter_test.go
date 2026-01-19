package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTAdapter(t *testing.T) {
	secret := "test-secret"
	expiration := 24

	adapter := NewJWTAdapter(secret, expiration)

	assert.NotNil(t, adapter)
	assert.Equal(t, secret, adapter.secret)
	assert.Equal(t, expiration, adapter.expiration)
}

func TestJWTAdapter_Generate_Success(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24
	adapter := NewJWTAdapter(secret, expiration)

	userID := int64(123)
	email := "test@example.com"

	token, err := adapter.Generate(userID, email)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Greater(t, len(token), 0)
}

func TestJWTAdapter_Generate_ContainsCorrectClaims(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24
	adapter := NewJWTAdapter(secret, expiration)

	userID := int64(123)
	email := "test@example.com"

	tokenString, err := adapter.Generate(userID, email)
	assert.NoError(t, err)

	// Parse the token to verify claims
	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	assert.NoError(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
}

func TestJWTAdapter_Validate_Success(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24
	adapter := NewJWTAdapter(secret, expiration)

	userID := int64(123)
	email := "test@example.com"

	// Generate a token first
	tokenString, err := adapter.Generate(userID, email)
	assert.NoError(t, err)

	// Validate the token
	claims, err := adapter.Validate(tokenString)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestJWTAdapter_Validate_InvalidSecret(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24
	adapter := NewJWTAdapter(secret, expiration)

	userID := int64(123)
	email := "test@example.com"

	// Generate a token with one secret
	tokenString, err := adapter.Generate(userID, email)
	assert.NoError(t, err)

	// Try to validate with different secret
	wrongAdapter := NewJWTAdapter("wrong-secret", expiration)
	claims, err := wrongAdapter.Validate(tokenString)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTAdapter_Validate_ExpiredToken(t *testing.T) {
	secret := "test-secret-key"
	expiration := -1 // Negative expiration means token is already expired
	adapter := NewJWTAdapter(secret, expiration)

	userID := int64(123)
	email := "test@example.com"

	// Generate a token that's already expired
	tokenString, err := adapter.Generate(userID, email)
	assert.NoError(t, err)

	// Wait a bit to ensure token is expired
	time.Sleep(100 * time.Millisecond)

	// Try to validate expired token
	claims, err := adapter.Validate(tokenString)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTAdapter_Validate_MalformedToken(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24
	adapter := NewJWTAdapter(secret, expiration)

	malformedToken := "not.a.valid.jwt.token"

	claims, err := adapter.Validate(malformedToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTAdapter_Validate_EmptyToken(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24
	adapter := NewJWTAdapter(secret, expiration)

	claims, err := adapter.Validate("")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTAdapter_Validate_InvalidSignature(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24
	adapter := NewJWTAdapter(secret, expiration)

	// Create a token with invalid signature by manually constructing it
	// This is a token with correct structure but wrong signature
	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.invalid-signature"

	claims, err := adapter.Validate(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTAdapter_RoundTrip(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24
	adapter := NewJWTAdapter(secret, expiration)

	userID := int64(456)
	email := "roundtrip@example.com"

	// Generate token
	tokenString, err := adapter.Generate(userID, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Validate token
	claims, err := adapter.Validate(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestJWTAdapter_Generate_DifferentUsers(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24
	adapter := NewJWTAdapter(secret, expiration)

	user1ID := int64(1)
	user1Email := "user1@example.com"
	user2ID := int64(2)
	user2Email := "user2@example.com"

	token1, err1 := adapter.Generate(user1ID, user1Email)
	token2, err2 := adapter.Generate(user2ID, user2Email)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, token1, token2) // Tokens should be different

	// Validate both tokens
	claims1, err1 := adapter.Validate(token1)
	claims2, err2 := adapter.Validate(token2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, user1ID, claims1.UserID)
	assert.Equal(t, user1Email, claims1.Email)
	assert.Equal(t, user2ID, claims2.UserID)
	assert.Equal(t, user2Email, claims2.Email)
}

func TestJWTAdapter_Validate_ImplementsInterface(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24
	adapter := NewJWTAdapter(secret, expiration)

	// Verify that JWTAdapter implements TokenGenerator and TokenValidator interfaces
	var _ domainuser.TokenGenerator = adapter
	var _ domainuser.TokenValidator = adapter
}

func TestJWTAdapter_Generate_ExpirationTime(t *testing.T) {
	secret := "test-secret-key"
	expiration := 2 // 2 hours
	adapter := NewJWTAdapter(secret, expiration)

	userID := int64(123)
	email := "test@example.com"

	tokenString, err := adapter.Generate(userID, email)
	assert.NoError(t, err)

	// Parse token to check expiration
	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	assert.NoError(t, err)
	assert.True(t, token.Valid)
	assert.NotNil(t, claims.ExpiresAt)

	// Check that expiration is approximately 2 hours from now
	expectedExpiration := time.Now().Add(2 * time.Hour)
	actualExpiration := claims.ExpiresAt.Time
	timeDiff := actualExpiration.Sub(expectedExpiration)

	// Allow 5 seconds difference for test execution time
	assert.True(t, timeDiff < 5*time.Second && timeDiff > -5*time.Second,
		"Expiration time should be approximately 2 hours from now")
}
