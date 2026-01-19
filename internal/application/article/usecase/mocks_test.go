package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
	"github.com/stretchr/testify/mock"
)

// mockArticleRepository is a mock implementation of Repository
type mockArticleRepository struct {
	mock.Mock
}

func (m *mockArticleRepository) Create(ctx context.Context, article *domainarticle.Article) (*domainarticle.Article, error) {
	args := m.Called(ctx, article)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainarticle.Article), args.Error(1)
}

func (m *mockArticleRepository) GetByID(ctx context.Context, id int64) (*domainarticle.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainarticle.Article), args.Error(1)
}

func (m *mockArticleRepository) Update(ctx context.Context, article *domainarticle.Article) (*domainarticle.Article, error) {
	args := m.Called(ctx, article)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainarticle.Article), args.Error(1)
}

func (m *mockArticleRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockArticleRepository) List(ctx context.Context, limit, offset int) ([]*domainarticle.Article, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainarticle.Article), args.Error(1)
}

func (m *mockArticleRepository) ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*domainarticle.Article, error) {
	args := m.Called(ctx, authorID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainarticle.Article), args.Error(1)
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

func (m *mockArticleCache) Get(ctx context.Context, id int64) (*domainarticle.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainarticle.Article), args.Error(1)
}

func (m *mockArticleCache) Set(ctx context.Context, id int64, article *domainarticle.Article) error {
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

func (m *mockArticleListCache) GetArticleList(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListArticlesResponse), args.Error(1)
}

func (m *mockArticleListCache) SetArticleList(ctx context.Context, limit, offset int, listResp *dto.ListArticlesResponse) error {
	args := m.Called(ctx, limit, offset, listResp)
	return args.Error(0)
}

func (m *mockArticleListCache) InvalidateArticleList(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

