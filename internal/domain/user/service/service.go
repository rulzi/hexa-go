package service

import "github.com/rulzi/hexa-go/internal/domain/user/port"

// Service provides domain-level business logic for users.
// It uses ports (interfaces) instead of concrete implementations.
type Service struct {
	repo           port.Repository
	tokenGen       port.TokenGenerator
	tokenValidator port.TokenValidator
	passwordHasher port.PasswordHasher
}

// NewService creates a new user service.
func NewService(
	repo port.Repository,
	tokenGen port.TokenGenerator,
	tokenValidator port.TokenValidator,
	passwordHasher port.PasswordHasher,
) *Service {
	return &Service{
		repo:           repo,
		tokenGen:       tokenGen,
		tokenValidator: tokenValidator,
		passwordHasher: passwordHasher,
	}
}
