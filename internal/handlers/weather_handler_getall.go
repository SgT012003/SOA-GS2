package handlers

import (
	"net/http"
	"strconv"

	"agri-ai-api/internal/models"

	"github.com/gin-gonic/gin"
)

// GetAllCachesHandler busca o histórico de consultas de clima
// @Summary Listar Cache Climático
// @Description Retorna todos os dados cacheados de clima com suporte a paginação
// @Tags weather
// @Accept json
// @Produce json
// @Param limit query int false "Limite de resultados (default 10)"
// @Param offset query int false "Deslocamento de resultados (default 0)"
// @Security BearerAuth
// @Success 200 {array} models.WeatherCache
// @Failure 500 {object} models.ProblemDetail "Erro interno"
// @Router /protected/weather/cache [get]
func (h *WeatherHandler) GetAllCachesHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	caches, err := h.weatherService.GetAllCaches(limit, offset)
	if err != nil {
		prob := models.NewTechnicalError("Erro Interno", err.Error(), c.Request.URL.Path, http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, prob)
		return
	}

	c.JSON(http.StatusOK, caches)
}
