package services

import (
	"context"

	"agri-ai-api/internal/dao"
	"agri-ai-api/internal/models"
)

// EngineService define as operações do motor
type EngineService interface {
	GenerateHarvestPlan(lat, lon float64) (*models.HarvestPlan, error)
	GenerateRiskAnalysis(ctx context.Context, lat, lon float64) (*models.RiskAnalysis, error)
	SelectIdealCrops(ctx context.Context, lat, lon float64, season string) ([]models.CropRecommendation, error)
}

type EngineServiceImpl struct {
	weatherService WeatherService
	cropDAO        dao.CropDAO
}

func NewEngineService(weatherService WeatherService, cropDAO dao.CropDAO) EngineService {
	return &EngineServiceImpl{
		weatherService: weatherService,
		cropDAO:        cropDAO,
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

// GenerateRiskAnalysis avalia riscos de geada ou incêndio baseados no clima
func (s *EngineServiceImpl) GenerateRiskAnalysis(ctx context.Context, lat, lon float64) (*models.RiskAnalysis, error) {
	weather, err := s.weatherService.GetWeather(lat, lon)
	if err != nil {
		return nil, err
	}

	temp := weather.Current.Temperature2m
	precip := weather.Current.Precipitation
	wind := weather.Current.WindSpeed10m

	riskScore := 0
	var factors []string

	// Risco de Geada (Temperaturas muito baixas)
	if temp < 4.0 {
		riskScore += 50
		factors = append(factors, "Alto Risco de Geada")
	} else if temp < 10.0 {
		riskScore += 20
		factors = append(factors, "Risco Moderado de Frio Extremo")
	}

	// Risco de Incêndio / Seca Severa (Temperatura alta + sem chuva + vento forte)
	if temp > 30.0 && precip < 1.0 {
		if wind > 15.0 {
			riskScore += 60
			factors = append(factors, "Risco Crítico de Incêndio Espontâneo")
		} else {
			riskScore += 30
			factors = append(factors, "Risco de Seca/Estresse Hídrico")
		}
	}

	// Limita score a 100
	if riskScore > 100 {
		riskScore = 100
	}

	if riskScore == 0 {
		factors = append(factors, "Nenhum risco iminente detectado")
	}

	return &models.RiskAnalysis{
		OverallRiskScore: riskScore,
		RiskFactors:      factors,
	}, nil
}

// SelectIdealCrops avalia qual cultura do banco se adapta melhor ao clima
func (s *EngineServiceImpl) SelectIdealCrops(ctx context.Context, lat, lon float64, season string) ([]models.CropRecommendation, error) {
	crops, err := s.cropDAO.GetAllCrops(ctx)
	if err != nil {
		return nil, err
	}

	weather, err := s.weatherService.GetWeather(lat, lon)
	if err != nil {
		return nil, err
	}

	temp := weather.Current.Temperature2m
	
	// Como a API Open-Meteo current retorna chuva do momento, 
	// para uma heurística real precisaríamos do histórico anual de chuvas.
	// Para este MVP, vamos simular uma "chuva mensal projetada" baseada na atual.
	projPrecip := weather.Current.Precipitation * 30 

	var recommendations []models.CropRecommendation

	for _, c := range crops {
		if season != "" && c.IdealSeason != season {
			continue // Pula se o usuário filtrou e a estação não bate
		}

		matchPercentage := 100.0
		var justification string

		// Penaliza por temperatura
		if temp < c.MinTemp {
			diff := c.MinTemp - temp
			matchPercentage -= diff * 5
			justification = "Temperatura atual muito baixa para esta cultura. "
		} else if temp > c.MaxTemp {
			diff := temp - c.MaxTemp
			matchPercentage -= diff * 5
			justification = "Temperatura atual muito alta para esta cultura. "
		}

		// Penaliza por precipitação projetada
		if projPrecip < c.MinPrecipitation {
			diff := c.MinPrecipitation - projPrecip
			matchPercentage -= (diff / c.MinPrecipitation) * 50
			justification += "Precipitação abaixo do ideal."
		} else if projPrecip > c.MaxPrecipitation {
			diff := projPrecip - c.MaxPrecipitation
			matchPercentage -= (diff / c.MaxPrecipitation) * 50
			justification += "Precipitação acima do ideal."
		}

		if matchPercentage < 0 {
			matchPercentage = 0
		}

		if justification == "" {
			justification = "Condições climáticas altamente compatíveis."
		}

		recommendations = append(recommendations, models.CropRecommendation{
			CropName:        c.Name,
			MatchPercentage: matchPercentage,
			Justification:   justification,
		})
	}

	return recommendations, nil
}
