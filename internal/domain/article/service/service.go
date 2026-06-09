package service

import "github.com/rulzi/hexa-go/internal/domain/article/port"

// Service provides domain-level business logic for articles.
type Service struct {
	repo port.Repository
}

// NewService creates a new article service.
func NewService(repo port.Repository) *Service {
	return &Service{repo: repo}
}
