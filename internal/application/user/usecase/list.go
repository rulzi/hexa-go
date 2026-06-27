package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	userport "github.com/rulzi/hexa-go/internal/domain/user/port"
)

type listUserDeps struct {
	userRepo userport.Repository
}

// ListUser handles user listing with pagination
type ListUser struct {
	deps listUserDeps
}

// NewListUser creates a new ListUser use case.
func NewListUser(userRepo userport.Repository) *ListUser {
	return &ListUser{deps: listUserDeps{userRepo: userRepo}}
}

// Execute lists users with pagination
func (uc *ListUser) Execute(ctx context.Context, limit, offset int) (*dto.ListUsersResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := uc.deps.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.deps.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

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
