package shared

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/grgn-stack/pkg/config"
)

// MockDatabase implements IDatabase for testing
type MockDatabase struct {
	pingError error
}

func (m *MockDatabase) Ping(ctx context.Context) error {
	return m.pingError
}

func (m *MockDatabase) Close(ctx context.Context) error {
	return nil
}

func (m *MockDatabase) VerifyConnectivity(ctx context.Context) error {
	return m.pingError
}

func (m *MockDatabase) ExecuteRead(ctx context.Context, work neo4j.ManagedTransactionWork) (any, error) {
	return nil, nil
}

func (m *MockDatabase) ExecuteWrite(ctx context.Context, work neo4j.ManagedTransactionWork) (any, error) {
	return nil, nil
}

func (m *MockDatabase) NewSession(ctx context.Context, config neo4j.SessionConfig) neo4j.SessionWithContext {
	return nil
}

func (m *MockDatabase) GetDriver() neo4j.DriverWithContext {
	return nil
}

func newTestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Environment: "test",
		},
		App: config.AppConfig{
			Version: "1.0.0-test",
		},
	}
}

func TestPingHandler_HandlePing_Healthy(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock database that returns no error (healthy)
	mockDB := &MockDatabase{pingError: nil}
	cfg := newTestConfig()
	handler := NewPingHandler(mockDB, cfg)

	// Create test router
	r := gin.Default()
	r.GET("/ping", handler.HandlePing)

	// Create test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"pong"`)
	assert.Contains(t, w.Body.String(), `"database":"healthy"`)
	assert.Contains(t, w.Body.String(), `"environment":"test"`)
	assert.Contains(t, w.Body.String(), `"version":"1.0.0-test"`)
	assert.NotContains(t, w.Body.String(), `"error"`)
}

func TestPingHandler_HandlePing_Unhealthy(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock database that returns an error (unhealthy)
	mockDB := &MockDatabase{pingError: errors.New("connection refused")}
	cfg := newTestConfig()
	handler := NewPingHandler(mockDB, cfg)

	// Create test router
	r := gin.Default()
	r.GET("/ping", handler.HandlePing)

	// Create test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"pong"`)
	assert.Contains(t, w.Body.String(), `"database":"unhealthy"`)
	assert.Contains(t, w.Body.String(), `"error":"connection refused"`)
}

func TestPingHandler_CheckHealth_Healthy(t *testing.T) {
	mockDB := &MockDatabase{pingError: nil}
	cfg := newTestConfig()
	handler := NewPingHandler(mockDB, cfg)

	response, err := handler.CheckHealth(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "pong", response.Message)
	assert.Equal(t, "healthy", response.Database)
	assert.Empty(t, response.Error)
}

func TestPingHandler_CheckHealth_Unhealthy(t *testing.T) {
	mockDB := &MockDatabase{pingError: errors.New("database unavailable")}
	cfg := newTestConfig()
	handler := NewPingHandler(mockDB, cfg)

	response, err := handler.CheckHealth(context.Background())

	assert.Error(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "pong", response.Message)
	assert.Equal(t, "unhealthy", response.Database)
	assert.Equal(t, "database unavailable", response.Error)
}
