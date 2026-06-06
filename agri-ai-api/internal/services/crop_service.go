package services

import (
	"context"

	"agri-ai-api/internal/dao"
	"agri-ai-api/internal/models"
)

// CropService define as operações para culturas
type CropService interface {
	GetAllCrops(ctx context.Context) ([]models.CropProfile, error)
}

type CropServiceImpl struct {
	cropDAO dao.CropDAO
}

func NewCropService(cropDAO dao.CropDAO) CropService {
	return &CropServiceImpl{
		cropDAO: cropDAO,
	}
}

// GetAllCrops busca todas as culturas
func (s *CropServiceImpl) GetAllCrops(ctx context.Context) ([]models.CropProfile, error) {
	return s.cropDAO.GetAllCrops(ctx)
}
