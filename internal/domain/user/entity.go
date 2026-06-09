package user

import "time"

// User represents the user entity in the domain
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hidden from JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
