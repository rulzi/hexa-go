package usecase

import (
	"context"

	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	"github.com/stretchr/testify/mock"
)

// mockArticleRepository is a mock implementation of Repository
type mockArticleRepository struct {
	mock.Mock
}

func (m *mockArticleRepository) Create(ctx context.Context, article *articleentity.Article) (*articleentity.Article, error) {
	args := m.Called(ctx, article)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*articleentity.Article), args.Error(1)
}

func (m *mockArticleRepository) GetByID(ctx context.Context, id int64) (*articleentity.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*articleentity.Article), args.Error(1)
}

func (m *mockArticleRepository) Update(ctx context.Context, article *articleentity.Article) (*articleentity.Article, error) {
	args := m.Called(ctx, article)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*articleentity.Article), args.Error(1)
}

func (m *mockArticleRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockArticleRepository) List(ctx context.Context, limit, offset int) ([]*articleentity.Article, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*articleentity.Article), args.Error(1)
}

func (m *mockArticleRepository) ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*articleentity.Article, error) {
	args := m.Called(ctx, authorID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*articleentity.Article), args.Error(1)
}

func (m *mockArticleRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockArticleRepository) CountByAuthor(ctx context.Context, authorID int64) (int64, error) {
	args := m.Called(ctx, authorID)
	return args.Get(0).(int64), args.Error(1)
}

// mockArticleCache is a mock implementation of Cache
type mockArticleCache struct {
	mock.Mock
}

func (m *mockArticleCache) Get(ctx context.Context, id int64) (*articleentity.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*articleentity.Article), args.Error(1)
}

func (m *mockArticleCache) Set(ctx context.Context, id int64, article *articleentity.Article) error {
	args := m.Called(ctx, id, article)
	return args.Error(0)
}

func (m *mockArticleCache) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockArticleCache) InvalidateList(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// mockArticleListCache is a mock implementation of ArticleListCache
type mockArticleListCache struct {
	mock.Mock
}

func (m *mockArticleListCache) GetList(ctx context.Context, limit, offset int) (*ArticleListPage, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ArticleListPage), args.Error(1)
}

func (m *mockArticleListCache) SetList(ctx context.Context, limit, offset int, page *ArticleListPage) error {
	args := m.Called(ctx, limit, offset, page)
	return args.Error(0)
}

func (m *mockArticleListCache) InvalidateList(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
