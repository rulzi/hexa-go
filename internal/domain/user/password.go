package user

// PasswordHasher is a port for password hashing operations
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) bool
}

