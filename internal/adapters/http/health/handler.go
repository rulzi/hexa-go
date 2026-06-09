package health

import (
	"context"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
)

const checkTimeout = 3 * time.Second

// Handler handles health check requests.
type Handler struct {
	db    *sql.DB
	redis *redis.Client
}

// NewHandler creates a new health check handler.
func NewHandler(db *sql.DB, redisClient *redis.Client) *Handler {
	return &Handler{
		db:    db,
		redis: redisClient,
	}
}

type checkResult struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Check handles GET /health and verifies MySQL and Redis connectivity.
func (h *Handler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), checkTimeout)
	defer cancel()

	checks := gin.H{}
	healthy := true

	mysqlResult := h.checkMySQL(ctx)
	checks["mysql"] = mysqlResult
	if mysqlResult.Status != "up" {
		healthy = false
	}

	redisResult := h.checkRedis(ctx)
	checks["redis"] = redisResult
	if redisResult.Status == "down" {
		healthy = false
	}

	data := gin.H{
		"status": "ok",
		"checks": checks,
	}
	if !healthy {
		data["status"] = "degraded"
		c.JSON(response.StatusCode.ServiceUnavailable(), response.StandardResponse{
			Status:  response.StatusError,
			Message: "Service is unhealthy",
			Data:    data,
		})
		return
	}

	response.SuccessResponseOK(c, "Service is healthy", data)
}

func (h *Handler) checkMySQL(ctx context.Context) checkResult {
	if h.db == nil {
		return checkResult{Status: "down", Error: "database not configured"}
	}

	if err := h.db.PingContext(ctx); err != nil {
		return checkResult{Status: "down", Error: err.Error()}
	}

	return checkResult{Status: "up"}
}

func (h *Handler) checkRedis(ctx context.Context) checkResult {
	if h.redis == nil {
		return checkResult{Status: "disabled"}
	}

	if err := h.redis.Ping(ctx).Err(); err != nil {
		return checkResult{Status: "down", Error: err.Error()}
	}

	return checkResult{Status: "up"}
}
