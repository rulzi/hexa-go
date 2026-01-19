package media

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockRepository is a mock implementation of Repository for testing
type mockRepository struct {
	createFunc  func(ctx context.Context, media *Media) (*Media, error)
	getByIDFunc func(ctx context.Context, id int64) (*Media, error)
	updateFunc  func(ctx context.Context, media *Media) (*Media, error)
	deleteFunc  func(ctx context.Context, id int64) error
	listFunc    func(ctx context.Context, limit, offset int) ([]*Media, error)
	countFunc   func(ctx context.Context) (int64, error)
}

func (m *mockRepository) Create(ctx context.Context, media *Media) (*Media, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, media)
	}
	return nil, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id int64) (*Media, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) Update(ctx context.Context, media *Media) (*Media, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, media)
	}
	return nil, nil
}

func (m *mockRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockRepository) List(ctx context.Context, limit, offset int) ([]*Media, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *mockRepository) Count(ctx context.Context) (int64, error) {
	if m.countFunc != nil {
		return m.countFunc(ctx)
	}
	return 0, nil
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name string
		repo Repository
	}{
		{
			name: "create service with repository",
			repo: &mockRepository{},
		},
		{
			name: "create service with nil repository",
			repo: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.repo)
			assert.NotNil(t, service)
			assert.Equal(t, tt.repo, service.repo)
		})
	}
}

