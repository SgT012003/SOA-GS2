package handlers

import (
	"net/http"
	"strconv"

	"agri-ai-api/internal/services"

	"github.com/gin-gonic/gin"
)

// HarvestEngineHandler analisa os dados e gera um plano de colheita
// @Summary Predição de Colheita
// @Description Avalia as condições climáticas e indica a viabilidade de colheita
// @Tags engine
// @Accept json
// @Produce json
// @Param lat query number true "Latitude"
// @Param lon query number true "Longitude"
// @Security BearerAuth
// @Success 200 {object} models.HarvestPlan
// @Failure 400 {object} map[string]string "Parâmetros inválidos"
// @Failure 429 {object} map[string]string "Muitas requisições (Rate Limit)"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /protected/engine/harvest [get]
func HarvestEngineHandler(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lon, errLon := strconv.ParseFloat(lonStr, 64)

	if errLat != nil || errLon != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros lat e lon devem ser números válidos"})
		return
	}

	plan, err := services.GenerateHarvestPlan(lat, lon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar o plano de colheita", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}
