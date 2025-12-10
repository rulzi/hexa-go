package article

// Service provides domain-level business logic for articles
type Service struct {
	repo Repository
}

// NewService creates a new article service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

