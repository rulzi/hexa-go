package user

// TokenClaims represents the claims extracted from a token
type TokenClaims struct {
	UserID int64
	Email  string
}

// TokenGenerator is a port for generating authentication tokens
type TokenGenerator interface {
	Generate(userID int64, email string) (string, error)
}

// TokenValidator is a port for validating authentication tokens
type TokenValidator interface {
	Validate(token string) (*TokenClaims, error)
}

