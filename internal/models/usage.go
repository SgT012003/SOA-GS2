package models

import "time"

// UsageLog representa um log de consumo da API
type UsageLog struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Endpoint  string    `json:"endpoint"`
	CreatedAt time.Time `json:"created_at"`
}
