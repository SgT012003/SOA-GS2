package models

import "time"

// OpenMeteoResponse representa a resposta básica da API do Open-Meteo
type OpenMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Current   struct {
		Time          string  `json:"time"`
		Temperature2m float64 `json:"temperature_2m"`
		Precipitation float64 `json:"precipitation"`
		WindSpeed10m  float64 `json:"wind_speed_10m"`
	} `json:"current"`
}

// WeatherCache representa um registro na tabela de cache
type WeatherCache struct {
	ID        int       `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	QueryDate time.Time `json:"query_date"`
	Data      string    `json:"data"` // JSON string da resposta
	CreatedAt time.Time `json:"created_at"`
}
