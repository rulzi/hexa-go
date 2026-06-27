package logging

import (
	"context"
	"testing"

	"github.com/rulzi/hexa-go/internal/adapters/contextkey"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	"github.com/stretchr/testify/require"
)

func TestNewRepoLogger(t *testing.T) {
	repoLog := NewRepoLogger(logger.NewSimpleLogger())
	require.NotNil(t, repoLog)
}

func TestRepoLogger_LogError_WithoutRequestID(t *testing.T) {
	repoLog := NewRepoLogger(logger.NewSimpleLogger())
	repoLog.LogError(context.Background(), "repo error: %s", "details")
}

func TestRepoLogger_LogError_WithRequestID(t *testing.T) {
	repoLog := NewRepoLogger(logger.NewSimpleLogger())
	ctx := contextkey.WithRequestID(context.Background(), "req-123")
	repoLog.LogError(ctx, "repo error: %s", "details")
}
