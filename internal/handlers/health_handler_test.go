package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"agri-ai-api/internal/dao"
	"agri-ai-api/internal/models"

	"github.com/gin-gonic/gin"
)

func TestHealthzHandler(t *testing.T) {
	// Configura Gin para o modo de teste
	gin.SetMode(gin.TestMode)

	// Garantir que não chame o banco nulo causando panic, se houver mock
	// No handler, se dao.DB for nil, ele retorna DEGRADED e código 503
	// Isso é bom, podemos testar essa lógica.
	dao.DB = nil

	r := gin.Default()
	r.GET("/healthz", HealthzHandler)

	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Como o banco não está inicializado neste escopo, esperamos 503 Service Unavailable
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", w.Code)
	}

	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.API != "DEGRADED" {
		t.Errorf("expected API status DEGRADED, got %s", response.API)
	}
	if response.Database != "DOWN (Nil pointer)" {
		t.Errorf("expected Database status DOWN (Nil pointer), got %s", response.Database)
	}
}
