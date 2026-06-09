package errmapper

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"not found", domainerrs.NewNotFound("USER_NOT_FOUND", "user not found"), http.StatusNotFound},
		{"validation", domainerrs.NewValidation("USER_NAME_REQUIRED", "name is required"), http.StatusBadRequest},
		{"conflict", domainerrs.NewConflict("USER_EMAIL_EXISTS", "email already exists"), http.StatusConflict},
		{"unauthorized", domainerrs.NewUnauthorized("USER_INVALID_CREDENTIALS", "invalid credentials"), http.StatusUnauthorized},
		{"unknown", errors.New("db connection refused"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, StatusCode(tt.err))
		})
	}
}

func TestRespond(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Respond(c, domainerrs.NewNotFound("USER_NOT_FOUND", "user not found"))

	assert.Equal(t, http.StatusNotFound, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "error", body["status"])
	assert.Equal(t, "user not found", body["message"])
}
