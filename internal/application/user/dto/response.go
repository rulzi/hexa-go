package dto

import "time"

// UserResponse represents the response DTO for user
type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListUsersResponse represents the response DTO for listing users
type ListUsersResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int64          `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

// LoginResponse represents the response DTO for login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
