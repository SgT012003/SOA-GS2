package handlers

import (
	"net/http"
	"strconv"

	"agri-ai-api/internal/models"
	"agri-ai-api/internal/services"

	"github.com/gin-gonic/gin"
)

type EngineHandler struct {
	engineService services.EngineService
}

func NewEngineHandler(engineService services.EngineService) *EngineHandler {
	return &EngineHandler{
		engineService: engineService,
	}
}

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
// @Failure 400 {object} models.ProblemDetail "Parâmetros inválidos"
// @Failure 429 {object} models.ProblemDetail "Muitas requisições (Rate Limit)"
// @Failure 500 {object} models.ProblemDetail "Erro interno"
// @Router /protected/engine/harvest [get]
func (h *EngineHandler) HarvestEngineHandler(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lon, errLon := strconv.ParseFloat(lonStr, 64)

	if errLat != nil || errLon != nil {
		prob := models.NewBusinessError("Parâmetros Inválidos", "Parâmetros lat e lon devem ser números válidos", c.Request.URL.Path, http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, prob)
		return
	}

	plan, err := h.engineService.GenerateHarvestPlan(lat, lon)
	if err != nil {
		prob := models.NewTechnicalError("Erro Interno", err.Error(), c.Request.URL.Path, http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, prob)
		return
	}

	c.JSON(http.StatusOK, plan)
}

// RiskAnalysisHandler realiza a análise de risco de incêndio ou geada
// @Summary Análise de Risco (Geada / Incêndio)
// @Description Avalia risco baseado em temperatura, vento e chuva
// @Tags engine
// @Accept json
// @Produce json
// @Param lat query number true "Latitude"
// @Param lon query number true "Longitude"
// @Security BearerAuth
// @Success 200 {object} models.RiskAnalysis
// @Failure 400 {object} models.ProblemDetail "Parâmetros inválidos"
// @Failure 500 {object} models.ProblemDetail "Erro interno"
// @Router /protected/engine/risk-analysis [get]
func (h *EngineHandler) RiskAnalysisHandler(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lon, errLon := strconv.ParseFloat(lonStr, 64)

	if errLat != nil || errLon != nil {
		prob := models.NewBusinessError("Parâmetros Inválidos", "Parâmetros lat e lon devem ser números válidos", c.Request.URL.Path, http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, prob)
		return
	}

	analysis, err := h.engineService.GenerateRiskAnalysis(c.Request.Context(), lat, lon)
	if err != nil {
		prob := models.NewTechnicalError("Erro ao Gerar Análise", err.Error(), c.Request.URL.Path, http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, prob)
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// CropSelectorHandler recomenda a melhor cultura agrícola para a região
// @Summary Seleção de Cultivo
// @Description Faz um match climático entre as safras do banco e o clima da região
// @Tags engine
// @Accept json
// @Produce json
// @Param lat query number true "Latitude"
// @Param lon query number true "Longitude"
// @Param season query string false "Estação do ano (ex: summer, winter)"
// @Security BearerAuth
// @Success 200 {array} models.CropRecommendation
// @Failure 400 {object} models.ProblemDetail "Parâmetros inválidos"
// @Failure 500 {object} models.ProblemDetail "Erro interno"
// @Router /protected/engine/crop-selector [get]
func (h *EngineHandler) CropSelectorHandler(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	season := c.Query("season")

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lon, errLon := strconv.ParseFloat(lonStr, 64)

	if errLat != nil || errLon != nil {
		prob := models.NewBusinessError("Parâmetros Inválidos", "Parâmetros lat e lon devem ser números válidos", c.Request.URL.Path, http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, prob)
		return
	}

	recommendations, err := h.engineService.SelectIdealCrops(c.Request.Context(), lat, lon, season)
	if err != nil {
		prob := models.NewTechnicalError("Erro ao Selecionar Cultura", err.Error(), c.Request.URL.Path, http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, prob)
		return
	}

	c.JSON(http.StatusOK, recommendations)
}
