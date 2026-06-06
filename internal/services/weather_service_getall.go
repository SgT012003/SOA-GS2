package services

import "agri-ai-api/internal/models"

// GetAllCaches retorna o histórico de caches salvos no banco com paginação
func (s *WeatherServiceImpl) GetAllCaches(limit, offset int) ([]models.WeatherCache, error) {
	return s.weatherDAO.GetAllCaches(limit, offset)
}
