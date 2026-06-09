package contextkey

import "context"

type requestIDKey struct{}

// RequestIDHeader is the HTTP header name used to propagate request IDs.
const RequestIDHeader = "X-Request-ID"

// WithRequestID returns a context with the given request ID.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

// RequestID returns the request ID from context, or empty string if not set.
func RequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}
