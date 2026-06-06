package models

// HealthResponse padroniza o retorno do healthcheck da API
type HealthResponse struct {
	API       string `json:"api" example:"UP"`
	Database  string `json:"database" example:"UP"`
	Uptime    string `json:"uptime" example:"14m3s"`
	Timestamp string `json:"timestamp" example:"2026-06-06T11:20:00Z"`
}
