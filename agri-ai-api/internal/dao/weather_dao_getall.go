package dao

import "agri-ai-api/internal/models"

// GetAllCaches retorna o histórico de caches salvos no banco com paginação
func (dao *WeatherDAOImpl) GetAllCaches(limit, offset int) ([]models.WeatherCache, error) {
	query := `SELECT id, latitude, longitude, query_date, data, created_at 
	          FROM weather_cache 
	          ORDER BY created_at DESC 
	          LIMIT $1 OFFSET $2`

	rows, err := dao.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var caches []models.WeatherCache
	for rows.Next() {
		var c models.WeatherCache
		if err := rows.Scan(&c.ID, &c.Latitude, &c.Longitude, &c.QueryDate, &c.Data, &c.CreatedAt); err != nil {
			return nil, err
		}
		caches = append(caches, c)
	}
	return caches, nil
}
