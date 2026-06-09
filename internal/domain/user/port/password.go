package port

//go:generate mockgen -destination=mocks/mock_password_hasher.go -package=mocks github.com/rulzi/hexa-go/internal/domain/user/port PasswordHasher

// PasswordHasher is a port for password hashing operations.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) bool
}
