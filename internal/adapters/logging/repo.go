package logging

import (
	"context"
	"log"

	"github.com/rulzi/hexa-go/internal/adapters/contextkey"
)

// LogRepoError logs repository-layer errors with request_id from context.
func LogRepoError(ctx context.Context, format string, args ...interface{}) {
	requestID := contextkey.RequestID(ctx)
	if requestID != "" {
		log.Printf("request_id=%s "+format, append([]interface{}{requestID}, args...)...)
		return
	}
	log.Printf(format, args...)
}
