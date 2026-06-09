package logging

import (
	"context"
	"fmt"

	"github.com/rulzi/hexa-go/internal/adapters/contextkey"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// RepoLogger logs repository-layer errors with structured logging.
type RepoLogger struct {
	logger logger.Logger
}

// NewRepoLogger creates a RepoLogger.
func NewRepoLogger(appLogger logger.Logger) *RepoLogger {
	return &RepoLogger{logger: appLogger}
}

// LogError logs repository-layer errors with request_id from context.
func (l *RepoLogger) LogError(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if requestID := contextkey.RequestID(ctx); requestID != "" {
		l.logger.ErrorWithFields(msg, map[string]interface{}{
			"request_id": requestID,
		})
		return
	}
	l.logger.Error(msg)
}
