package usecase

import (
	"context"

	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
)

// NoopCache is a null-object implementation used when Redis is unavailable.
type NoopCache struct{}

var (
	_ ArticleCache     = NoopCache{}
	_ ArticleListCache = NoopCache{}
)

func (NoopCache) Get(_ context.Context, _ int64) (*articleentity.Article, error) {
	return nil, nil
}

func (NoopCache) Set(_ context.Context, _ int64, _ *articleentity.Article) error {
	return nil
}

func (NoopCache) Delete(_ context.Context, _ int64) error {
	return nil
}

func (NoopCache) InvalidateList(_ context.Context) error {
	return nil
}

func (NoopCache) GetList(_ context.Context, _, _ int) (*ArticleListPage, error) {
	return nil, nil
}

func (NoopCache) SetList(_ context.Context, _, _ int, _ *ArticleListPage) error {
	return nil
}
