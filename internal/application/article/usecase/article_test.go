package usecase

import (
	"testing"

	articleservice "github.com/rulzi/hexa-go/internal/domain/article/service"
	"github.com/stretchr/testify/assert"
)

func TestNewArticleUseCase(t *testing.T) {
	repo := &mockArticleRepository{}
	service := articleservice.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	assert.NotNil(t, uc)
	assert.NotNil(t, uc.Create)
	assert.NotNil(t, uc.Get)
	assert.NotNil(t, uc.List)
	assert.NotNil(t, uc.Update)
	assert.NotNil(t, uc.Delete)
}
