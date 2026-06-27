package entity

import "time"

const minPasswordLength = 6

// User represents the user entity in the domain
type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ValidatePassword checks plain-text password business rules before hashing.
func ValidatePassword(password string) error {
	if password == "" {
		return NewPasswordRequired()
	}
	if len(password) < minPasswordLength {
		return NewPasswordTooShort()
	}
	return nil
}

// ValidateRegistration checks user input before persistence.
func ValidateRegistration(name, email, plainPassword string) error {
	if err := ValidatePassword(plainPassword); err != nil {
		return err
	}
	if name == "" {
		return NewNameRequired()
	}
	if email == "" {
		return NewEmailRequired()
	}
	return nil
}

// Validate validates the user entity
func (u *User) Validate() error {
	if u.Name == "" {
		return NewNameRequired()
	}
	if u.Email == "" {
		return NewEmailRequired()
	}
	if u.Password == "" {
		return NewPasswordRequired()
	}
	return nil
}
