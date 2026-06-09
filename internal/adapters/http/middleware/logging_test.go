package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

type captureLogger struct {
	lastMsg    string
	lastFields map[string]interface{}
}

func (l *captureLogger) Debug(string)                                         {}
func (l *captureLogger) DebugWithFields(string, map[string]interface{})       {}
func (l *captureLogger) Info(string)                                        {}
func (l *captureLogger) InfoWithFields(msg string, fields map[string]interface{}) {
	l.lastMsg = msg
	l.lastFields = fields
}
func (l *captureLogger) Warn(string)                                          {}
func (l *captureLogger) WarnWithFields(msg string, fields map[string]interface{}) {
	l.lastMsg = msg
	l.lastFields = fields
}
func (l *captureLogger) Error(string)                                         {}
func (l *captureLogger) ErrorWithFields(msg string, fields map[string]interface{}) {
	l.lastMsg = msg
	l.lastFields = fields
}
func (l *captureLogger) Fatal(string)                                         {}
func (l *captureLogger) FatalWithFields(string, map[string]interface{})       {}
func (l *captureLogger) WithFields(fields map[string]interface{}) logger.Logger {
	return l
}

func TestLoggingMiddleware_LogsStructuredFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	appLogger := &captureLogger{}

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.Use(LoggingMiddleware(appLogger))
	router.GET("/items", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "http request completed", appLogger.lastMsg)
	assert.Equal(t, http.MethodGet, appLogger.lastFields["method"])
	assert.Equal(t, "/items", appLogger.lastFields["path"])
	assert.Equal(t, http.StatusOK, appLogger.lastFields["status"])
	assert.NotEmpty(t, appLogger.lastFields["request_id"])
	assert.Contains(t, appLogger.lastFields, "latency_ms")
}
