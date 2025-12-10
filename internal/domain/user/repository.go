package user

import "context"

// Repository is the driven port (interface) for user persistence
// This defines what the domain needs, not how it's implemented
type Repository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) (*User, error)

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (*User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *User) (*User, error)

	// Delete deletes a user by ID
	Delete(ctx context.Context, id int64) error

	// List retrieves all users with pagination
	List(ctx context.Context, limit, offset int) ([]*User, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
}

