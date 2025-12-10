package user

import (
	"context"

	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// ListUsersUseCase handles listing users with pagination
type ListUsersUseCase struct {
	userRepo domainuser.Repository
}

// NewListUsersUseCase creates a new ListUsersUseCase
func NewListUsersUseCase(userRepo domainuser.Repository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo: userRepo,
	}
}

// Execute executes the list users use case
func (uc *ListUsersUseCase) Execute(ctx context.Context, limit, offset int) (*ListUsersResponse, error) {
	// Default pagination
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Get users
	users, err := uc.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	userResponses := make([]UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}

	return &ListUsersResponse{
		Users:  userResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
