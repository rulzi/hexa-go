package entity

import (
	"fmt"
	"testing"

	domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"
	"github.com/stretchr/testify/assert"
)

func TestMediaErrors(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		check func(error) bool
		kind  func(error) bool
	}{
		{"not found", NewMediaNotFound(), IsMediaNotFound, domainerrs.IsNotFound},
		{"name required", NewNameRequired(), IsNameRequired, domainerrs.IsValidation},
		{"path required", NewPathRequired(), IsPathRequired, domainerrs.IsValidation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.check(tt.err))
			assert.True(t, tt.kind(tt.err))
		})
	}
}

func TestMediaErrors_Wrapped(t *testing.T) {
	wrapped := fmt.Errorf("usecase: %w", NewMediaNotFound())
	assert.True(t, IsMediaNotFound(wrapped))
	assert.True(t, domainerrs.IsNotFound(wrapped))
}
