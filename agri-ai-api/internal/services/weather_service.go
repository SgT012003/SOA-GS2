package services

import (
	"encoding/json"
	"log"
	"time"

	"agri-ai-api/internal/dao"
	"agri-ai-api/internal/models"
)

// WeatherService define as operações de clima
type WeatherService interface {
	GetWeather(lat, lon float64) (*models.OpenMeteoResponse, error)
	GetAllCaches(limit, offset int) ([]models.WeatherCache, error)
}

type WeatherServiceImpl struct {
	weatherDAO dao.WeatherDAO
}

func NewWeatherService(weatherDAO dao.WeatherDAO) WeatherService {
	return &WeatherServiceImpl{
		weatherDAO: weatherDAO,
	}
}

// GetWeather coordena a busca de clima, verificando primeiro no cache
func (s *WeatherServiceImpl) GetWeather(lat, lon float64) (*models.OpenMeteoResponse, error) {
	now := time.Now()

	cached, err := s.weatherDAO.GetCache(lat, lon, now)
	if err != nil {
		log.Printf("Erro ao buscar do cache: %v", err)
	} else if cached != nil {
		log.Println("Retornando dados do banco de dados (Cache hit)")
		var response models.OpenMeteoResponse
		if err := json.Unmarshal([]byte(cached.Data), &response); err == nil {
			return &response, nil
		}
	}

	log.Println("Buscando dados na Open-Meteo (Cache miss)")
	data, err := FetchOpenMeteo(lat, lon)
	if err != nil {
		return nil, err
	}

	var response models.OpenMeteoResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	go func() {
		if err := s.weatherDAO.SaveCache(lat, lon, now, string(data)); err != nil {
			log.Printf("Falha ao salvar cache no banco: %v", err)
		}
	}()

	return &response, nil
}
