package handler

import (
	"net/http"
	"time"

	"github.com/ai-fitness-planner/backend/internal/pkg/database"
	"github.com/ai-fitness-planner/backend/internal/pkg/redis"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp int64             `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// HealthCheck handles GET /health
// @Summary Health check
// @Description Check the health status of the API and its dependencies
// @Tags System
// @Produce json
// @Success 200 {object} HealthResponse "Service is healthy"
// @Failure 503 {object} HealthResponse "Service is unhealthy"
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	services := make(map[string]string)

	// Check database
	if database.DB != nil {
		sqlDB, err := database.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			services["database"] = "unhealthy"
		} else {
			services["database"] = "healthy"
		}
	} else {
		services["database"] = "not_initialized"
	}

	// Check Redis
	if redis.Rdb != nil {
		if err := redis.Rdb.Ping(c.Request.Context()).Err(); err != nil {
			services["redis"] = "unhealthy"
		} else {
			services["redis"] = "healthy"
		}
	} else {
		services["redis"] = "not_initialized"
	}

	// Determine overall status
	status := "healthy"
	for _, serviceStatus := range services {
		if serviceStatus != "healthy" {
			status = "unhealthy"
			break
		}
	}

	httpStatus := http.StatusOK
	if status == "unhealthy" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, HealthResponse{
		Status:    status,
		Timestamp: time.Now().Unix(),
		Services:  services,
	})
}
