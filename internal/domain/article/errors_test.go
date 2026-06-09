package article

import (
	"fmt"
	"testing"

	domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"
	"github.com/stretchr/testify/assert"
)

func TestArticleErrors(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		check func(error) bool
		kind  func(error) bool
	}{
		{"not found", NewArticleNotFound(), IsArticleNotFound, domainerrs.IsNotFound},
		{"title required", NewTitleRequired(), IsTitleRequired, domainerrs.IsValidation},
		{"content required", NewContentRequired(), IsContentRequired, domainerrs.IsValidation},
		{"author id required", NewAuthorIDRequired(), IsAuthorIDRequired, domainerrs.IsValidation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.check(tt.err))
			assert.True(t, tt.kind(tt.err))
		})
	}
}

func TestArticleErrors_Wrapped(t *testing.T) {
	wrapped := fmt.Errorf("usecase: %w", NewArticleNotFound())
	assert.True(t, IsArticleNotFound(wrapped))
	assert.True(t, domainerrs.IsNotFound(wrapped))
}
