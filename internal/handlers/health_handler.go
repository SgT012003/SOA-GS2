package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"agri-ai-api/internal/dao"
	"agri-ai-api/internal/models"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// HealthzHandler retorna a saúde atual do sistema
// @Summary Liveness/Readiness Probe
// @Description Verifica status da API e conexões com o banco de dados
// @Tags system
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Failure 503 {object} models.HealthResponse
// @Router /healthz [get]
func HealthzHandler(c *gin.Context) {
	uptime := time.Since(startTime).String()

	status := "UP"
	dbStatus := "UP"

	// Tenta pingar o banco de dados
	if dao.DB == nil {
		dbStatus = "DOWN (Nil pointer)"
		status = "DEGRADED"
	} else if err := dao.DB.Ping(); err != nil {
		dbStatus = "DOWN (" + err.Error() + ")"
		status = "DEGRADED"
		slog.Error("HealthCheck falhou na conexão com DB", slog.String("error", err.Error()))
	}

	response := models.HealthResponse{
		API:       status,
		Database:  dbStatus,
		Uptime:    uptime,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if status == "DEGRADED" {
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
