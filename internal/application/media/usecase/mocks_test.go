package usecase

import (
	"context"
	"io"

	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
	"github.com/stretchr/testify/mock"
)

// mockMediaRepository is a mock implementation of Repository
type mockMediaRepository struct {
	mock.Mock
}

func (m *mockMediaRepository) Create(ctx context.Context, media *domainmedia.Media) (*domainmedia.Media, error) {
	args := m.Called(ctx, media)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainmedia.Media), args.Error(1)
}

func (m *mockMediaRepository) GetByID(ctx context.Context, id int64) (*domainmedia.Media, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainmedia.Media), args.Error(1)
}

func (m *mockMediaRepository) Update(ctx context.Context, media *domainmedia.Media) (*domainmedia.Media, error) {
	args := m.Called(ctx, media)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainmedia.Media), args.Error(1)
}

func (m *mockMediaRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockMediaRepository) List(ctx context.Context, limit, offset int) ([]*domainmedia.Media, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainmedia.Media), args.Error(1)
}

func (m *mockMediaRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// mockMediaStorage is a mock implementation of Storage
type mockMediaStorage struct {
	mock.Mock
}

func (m *mockMediaStorage) Save(ctx context.Context, filename string, file io.Reader) (string, error) {
	args := m.Called(ctx, filename, file)
	return args.String(0), args.Error(1)
}

func (m *mockMediaStorage) Delete(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *mockMediaStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
