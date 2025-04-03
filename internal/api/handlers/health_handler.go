package handlers

import (
	"log"
	"time"

	"beautyessentials.com/internal/api/responses"
	"beautyessentials.com/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	healthService interfaces.HealthService
	respHelper    *responses.ResponseHelper
}

// NewHealthHandler creates a new instance of HealthHandler
func NewHealthHandler(
	healthService interfaces.HealthService,
	respHelper *responses.ResponseHelper,
) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
		respHelper:    respHelper,
	}
}

// HealthCheck handles the health check endpoint
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	log.Println("Health check endpoint called")
	
	// Call the service layer
	healthy, err := h.healthService.CheckHealth(c)
	if err != nil {
		h.respHelper.SendError(c, "Health check failed", err.Error(), responses.HTTPUnprocessableEntity)
		return
	}
	
	status := "ok"
	if !healthy {
		status = "unhealthy"
	}
	
	h.respHelper.OkResponse(c, gin.H{
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
	}, "Health check completed")
}