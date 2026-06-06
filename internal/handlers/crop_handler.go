package handlers

import (
	"net/http"

	"agri-ai-api/internal/models"
	"agri-ai-api/internal/services"

	"github.com/gin-gonic/gin"
)

type CropHandler struct {
	cropService services.CropService
}

func NewCropHandler(cropService services.CropService) *CropHandler {
	return &CropHandler{
		cropService: cropService,
	}
}

// GetAllCropsHandler retorna a lista de culturas
// @Summary Listar Culturas Agrícolas
// @Description Retorna todos os perfis de culturas (Crops) disponíveis no banco
// @Tags crops
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.CropProfile
// @Failure 500 {object} models.ProblemDetail "Erro interno"
// @Router /protected/crops [get]
func (h *CropHandler) GetAllCropsHandler(c *gin.Context) {
	crops, err := h.cropService.GetAllCrops(c.Request.Context())
	if err != nil {
		prob := models.NewTechnicalError("Erro Interno", err.Error(), c.Request.URL.Path, http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, prob)
		return
	}

	c.JSON(http.StatusOK, crops)
}
