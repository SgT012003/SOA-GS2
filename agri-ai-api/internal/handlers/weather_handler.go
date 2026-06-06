package handlers

import (
	"net/http"
	"strconv"

	"agri-ai-api/internal/services"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	weatherService services.WeatherService
}

func NewWeatherHandler(weatherService services.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

// GetWeatherHandler busca os dados de clima
// @Summary Buscar Clima Atual
// @Description Retorna a temperatura e precipitação de uma coordenada. Utiliza cache no banco de dados e Open-Meteo como fallback.
// @Tags weather
// @Accept json
// @Produce json
// @Param lat query number true "Latitude"
// @Param lon query number true "Longitude"
// @Security BearerAuth
// @Success 200 {object} models.OpenMeteoResponse
// @Failure 400 {object} map[string]string "Parâmetros inválidos"
// @Failure 500 {object} map[string]string "Erro interno ao buscar clima"
// @Router /protected/weather [get]
func (h *WeatherHandler) GetWeatherHandler(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lon, errLon := strconv.ParseFloat(lonStr, 64)

	if errLat != nil || errLon != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros lat e lon devem ser números válidos"})
		return
	}

	response, err := h.weatherService.GetWeather(lat, lon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar dados climáticos", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
