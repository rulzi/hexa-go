package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		limitStr   string
		offsetStr  string
		wantLimit  int
		wantOffset int
	}{
		{"defaults on empty", "", "", DefaultLimit, 0},
		{"defaults on invalid", "abc", "xyz", DefaultLimit, 0},
		{"valid values", "20", "5", 20, 5},
		{"clamps limit to max", "1000000", "0", MaxLimit, 0},
		{"rejects zero limit", "0", "0", DefaultLimit, 0},
		{"rejects negative limit", "-5", "0", DefaultLimit, 0},
		{"rejects negative offset", "10", "-1", 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit, offset := Parse(tt.limitStr, tt.offsetStr)
			assert.Equal(t, tt.wantLimit, limit)
			assert.Equal(t, tt.wantOffset, offset)
		})
	}
}
