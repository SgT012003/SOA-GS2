package services

import (
	"agri-ai-api/internal/models"
)

// GenerateHarvestPlan cria uma recomendação de colheita baseada no clima
func GenerateHarvestPlan(lat, lon float64) (*models.HarvestPlan, error) {
	// 1. Busca os dados climáticos atuais (reutilizando o serviço construído no Sprint 3)
	weather, err := GetWeather(lat, lon)
	if err != nil {
		return nil, err
	}

	temp := weather.Current.Temperature2m
	precip := weather.Current.Precipitation

	// 2. Motor de Regras Estatísticas/Preditivas Simuladas
	// Cenário Ideal: Temperatura amena e sem chuva
	score := 100
	recommendation := "Condições excelentes para colheita."

	// Penalização por chuva (colheita com chuva não é ideal)
	if precip > 0 {
		if precip < 2.0 {
			score -= 30
			recommendation = "Atenção: Chuva leve. Possível realizar colheita, mas com cautela."
		} else {
			score -= 80
			recommendation = "Não recomendado: Chuva forte impossibilita a colheita segura."
		}
	}

	// Penalização por temperatura extrema
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

	// Limita o score entre 0 e 100
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
