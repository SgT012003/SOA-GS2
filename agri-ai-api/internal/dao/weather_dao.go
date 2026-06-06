package dao

import (
	"database/sql"
	"errors"
	"time"

	"agri-ai-api/internal/models"
)

// WeatherDAO define os métodos de acesso aos dados de cache climático
type WeatherDAO interface {
	GetCache(lat, lon float64, date time.Time) (*models.WeatherCache, error)
	SaveCache(lat, lon float64, date time.Time, data string) error
	GetAllCaches(limit, offset int) ([]models.WeatherCache, error)
}

// WeatherDAOImpl é a implementação
type WeatherDAOImpl struct {
	db *sql.DB
}

// NewWeatherDAO retorna a implementação
func NewWeatherDAO() *WeatherDAOImpl {
	return &WeatherDAOImpl{db: DB}
}

// GetCache busca dados cacheados por latitude, longitude e data
func (dao *WeatherDAOImpl) GetCache(lat, lon float64, date time.Time) (*models.WeatherCache, error) {
	cache := &models.WeatherCache{}
	query := `SELECT id, latitude, longitude, query_date, data, created_at 
	          FROM weather_cache 
	          WHERE latitude = $1 AND longitude = $2 AND query_date = $3`

	// Trunca a data para garantir precisão apenas no dia
	truncatedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	err := dao.db.QueryRow(query, lat, lon, truncatedDate).Scan(
		&cache.ID, &cache.Latitude, &cache.Longitude, &cache.QueryDate, &cache.Data, &cache.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	return cache, nil
}

// SaveCache armazena a resposta bruta no banco
func (dao *WeatherDAOImpl) SaveCache(lat, lon float64, date time.Time, data string) error {
	query := `INSERT INTO weather_cache (latitude, longitude, query_date, data) 
	          VALUES ($1, $2, $3, $4) 
	          ON CONFLICT (latitude, longitude, query_date) DO NOTHING`

	truncatedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	_, err := dao.db.Exec(query, lat, lon, truncatedDate, data)
	return err
}
