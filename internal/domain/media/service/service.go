package service

import "github.com/rulzi/hexa-go/internal/domain/media/port"

// Service provides domain-level business logic for media.
type Service struct {
	repo port.Repository
}

// NewService creates a new media service.
func NewService(repo port.Repository) *Service {
	return &Service{repo: repo}
}
