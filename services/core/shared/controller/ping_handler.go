package shared

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/grgn-stack/pkg/config"
)

// PingHandler handles health check requests for the application.
// It checks database connectivity and returns the health status.
type PingHandler struct {
	db     IDatabase
	config *config.Config
}

// PingResponse represents the response from the ping endpoint.
type PingResponse struct {
	Message     string `json:"message"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
	Database    string `json:"database"`
	Error       string `json:"error,omitempty"`
}

// NewPingHandler creates a new PingHandler with the given dependencies.
func NewPingHandler(db IDatabase, cfg *config.Config) *PingHandler {
	return &PingHandler{
		db:     db,
		config: cfg,
	}
}

// HandlePing processes the health check request.
// It verifies database connectivity and returns the service health status.
func (h *PingHandler) HandlePing(c *gin.Context) {
	response := PingResponse{
		Message:     "pong",
		Environment: h.config.Server.Environment,
		Version:     h.config.App.Version,
		Database:    "healthy",
	}

	// Check database connectivity with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		response.Database = "unhealthy"
		response.Error = err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CheckHealth performs a health check and returns the result.
// This method can be called programmatically without HTTP context.
func (h *PingHandler) CheckHealth(ctx context.Context) (*PingResponse, error) {
	response := &PingResponse{
		Message:     "pong",
		Environment: h.config.Server.Environment,
		Version:     h.config.App.Version,
		Database:    "healthy",
	}

	// Check database connectivity
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := h.db.Ping(checkCtx); err != nil {
		response.Database = "unhealthy"
		response.Error = err.Error()
		return response, err
	}

	return response, nil
}
