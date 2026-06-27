package auth

import (
	"context"

	"github.com/rulzi/hexa-go/internal/domain/user/entity"
	userport "github.com/rulzi/hexa-go/internal/domain/user/port"
)

// Result holds the outcome of a successful authentication.
type Result struct {
	User  *entity.User
	Token string
}

// Authenticate verifies credentials and issues an auth token.
func Authenticate(
	ctx context.Context,
	email, password string,
	repo userport.Repository,
	hasher userport.PasswordHasher,
	tokenGen userport.TokenGenerator,
) (*Result, error) {
	user, err := repo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, entity.NewInvalidCredentials()
	}

	if !hasher.Verify(user.Password, password) {
		return nil, entity.NewInvalidCredentials()
	}

	token, err := tokenGen.Generate(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &Result{User: user, Token: token}, nil
}
