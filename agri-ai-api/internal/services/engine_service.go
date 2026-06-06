package services

import (
	"agri-ai-api/internal/models"
)

// EngineService define as operações do motor
type EngineService interface {
	GenerateHarvestPlan(lat, lon float64) (*models.HarvestPlan, error)
}

type EngineServiceImpl struct {
	weatherService WeatherService
}

func NewEngineService(weatherService WeatherService) EngineService {
	return &EngineServiceImpl{
		weatherService: weatherService,
	}
}

// GenerateHarvestPlan cria uma recomendação de colheita baseada no clima
func (s *EngineServiceImpl) GenerateHarvestPlan(lat, lon float64) (*models.HarvestPlan, error) {
	weather, err := s.weatherService.GetWeather(lat, lon)
	if err != nil {
		return nil, err
	}

	temp := weather.Current.Temperature2m
	precip := weather.Current.Precipitation

	score := 100
	recommendation := "Condições excelentes para colheita."

	if precip > 0 {
		if precip < 2.0 {
			score -= 30
			recommendation = "Atenção: Chuva leve. Possível realizar colheita, mas com cautela."
		} else {
			score -= 80
			recommendation = "Não recomendado: Chuva forte impossibilita a colheita segura."
		}
	}

	if temp > 35.0 {
		score -= 40
		if score > 0 {
			recommendation = "Atenção: Temperatura extrema. Risco de degradação rápida dos frutos."
		}
	} else if temp < 5.0 {
		score -= 50
		if score > 0 {
			recommendation = "Atenção: Muito frio. Risco de geada ou danos às plantas."
		}
	}

	if score < 0 {
		score = 0
	}

	plan := &models.HarvestPlan{
		ViabilityScore: score,
		Recommendation: recommendation,
		Temperature:    temp,
		Precipitation:  precip,
	}

	return plan, nil
}
