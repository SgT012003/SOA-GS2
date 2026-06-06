package services

import (
	"encoding/json"
	"log"
	"time"

	"agri-ai-api/internal/dao"
	"agri-ai-api/internal/models"
)

var weatherDAO = dao.NewWeatherDAO()

// GetWeather coordena a busca de clima, verificando primeiro no cache
func GetWeather(lat, lon float64) (*models.OpenMeteoResponse, error) {
	now := time.Now()

	// 1. Tentar pegar do cache do banco
	cached, err := weatherDAO.GetCache(lat, lon, now)
	if err != nil {
		log.Printf("Erro ao buscar do cache: %v", err)
		// Continua para tentar da API externa, mesmo com erro no banco
	} else if cached != nil {
		log.Println("Retornando dados do banco de dados (Cache hit)")
		var response models.OpenMeteoResponse
		if err := json.Unmarshal([]byte(cached.Data), &response); err == nil {
			return &response, nil
		}
	}

	// 2. Se não tem no cache, buscar na API externa
	log.Println("Buscando dados na Open-Meteo (Cache miss)")
	data, err := FetchOpenMeteo(lat, lon)
	if err != nil {
		return nil, err
	}

	// Tentar fazer parse
	var response models.OpenMeteoResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	// 3. Salvar no cache
	// É feito de forma assíncrona para não bloquear o retorno
	go func() {
		if err := weatherDAO.SaveCache(lat, lon, now, string(data)); err != nil {
			log.Printf("Falha ao salvar cache no banco: %v", err)
		}
	}()

	return &response, nil
}
