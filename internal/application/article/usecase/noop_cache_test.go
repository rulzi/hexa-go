package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopCache_ImplementsInterfaces(t *testing.T) {
	var _ ArticleCache = NoopCache{}
	var _ ArticleListCache = NoopCache{}
}

func TestNoopCache_ArticleOperations(t *testing.T) {
	cache := NoopCache{}

	got, err := cache.Get(t.Context(), 1)
	assert.NoError(t, err)
	assert.Nil(t, got)
	assert.NoError(t, cache.Set(t.Context(), 1, nil))
	assert.NoError(t, cache.Delete(t.Context(), 1))
	assert.NoError(t, cache.InvalidateList(t.Context()))
}

func TestNoopCache_ListOperations(t *testing.T) {
	cache := NoopCache{}

	got, err := cache.GetList(t.Context(), 10, 0)
	assert.NoError(t, err)
	assert.Nil(t, got)
	assert.NoError(t, cache.SetList(t.Context(), 10, 0, nil))
}
