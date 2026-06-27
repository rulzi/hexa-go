package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	httparticle "github.com/rulzi/hexa-go/internal/adapters/http/article"
	httphealth "github.com/rulzi/hexa-go/internal/adapters/http/health"
	httpmedia "github.com/rulzi/hexa-go/internal/adapters/http/media"
	httpuser "github.com/rulzi/hexa-go/internal/adapters/http/user"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	log := logger.NewSimpleLogger()
	healthHandler := httphealth.NewHandler(nil, nil)

	router := NewRouter(
		httpuser.NewHandler(httpuser.Deps{}, log),
		httparticle.NewHandler(httparticle.Deps{}, log),
		httpmedia.NewHandler(httpmedia.Deps{}, log),
		healthHandler,
		nil,
		t.TempDir(),
		[]string{"http://localhost:3000"},
		log,
	)

	assert.NotNil(t, router)
}

func TestRouter_SetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logger.NewSimpleLogger()
	healthHandler := httphealth.NewHandler(nil, nil)

	router := NewRouter(
		httpuser.NewHandler(httpuser.Deps{}, log),
		httparticle.NewHandler(httparticle.Deps{}, log),
		httpmedia.NewHandler(httpmedia.Deps{}, log),
		healthHandler,
		nil,
		t.TempDir(),
		nil,
		log,
	)

	engine := gin.New()
	router.SetupRoutes(engine)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.NotEqual(t, http.StatusNotFound, rec.Code)
}
