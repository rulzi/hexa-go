package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
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
func (uc *ListUsersUseCase) Execute(ctx context.Context, limit, offset int) (*dto.ListUsersResponse, error) {
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
	userResponses := make([]dto.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = dto.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}

	return &dto.ListUsersResponse{
		Users:  userResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
