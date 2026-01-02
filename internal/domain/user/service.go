package user

// Service provides domain-level business logic for users
// It uses ports (interfaces) instead of concrete implementations
type Service struct {
	repo          Repository
	tokenGen      TokenGenerator
	tokenValidator TokenValidator
	passwordHasher PasswordHasher
}

// NewService creates a new user service
func NewService(repo Repository, tokenGen TokenGenerator, tokenValidator TokenValidator, passwordHasher PasswordHasher) *Service {
	return &Service{
		repo:           repo,
		tokenGen:       tokenGen,
		tokenValidator: tokenValidator,
		passwordHasher:  passwordHasher,
	}
}
