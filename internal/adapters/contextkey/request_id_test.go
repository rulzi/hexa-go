package contextkey

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithRequestIDAndRequestID(t *testing.T) {
	ctx := WithRequestID(context.Background(), "abc-123")
	assert.Equal(t, "abc-123", RequestID(ctx))
}

func TestRequestID_EmptyWhenNotSet(t *testing.T) {
	assert.Equal(t, "", RequestID(context.Background()))
	assert.Equal(t, "", RequestID(nil))
}
