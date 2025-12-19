package media

// Service provides domain-level business logic for media
type Service struct {
	repo Repository
}

// NewService creates a new media service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
