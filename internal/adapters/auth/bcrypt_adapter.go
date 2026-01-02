package auth

import "golang.org/x/crypto/bcrypt"

// BcryptPasswordHasher implements PasswordHasher using bcrypt
type BcryptPasswordHasher struct{}

// NewBcryptPasswordHasher creates a new bcrypt password hasher
func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{}
}

// Hash implements PasswordHasher interface
func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify implements PasswordHasher interface
func (h *BcryptPasswordHasher) Verify(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

