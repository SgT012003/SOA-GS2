package dao

import (
	"context"
	"database/sql"

	"agri-ai-api/internal/models"
)

// CropDAO define a interface de acesso a dados para Culturas Agrícolas
type CropDAO interface {
	GetAllCrops(ctx context.Context) ([]models.CropProfile, error)
}

type cropDAOImpl struct {
	db *sql.DB
}

// NewCropDAO cria a implementação baseada no banco de dados
func NewCropDAO(db *sql.DB) CropDAO {
	return &cropDAOImpl{db: db}
}

// GetAllCrops retorna todas as culturas do banco usando SQL puro
func (dao *cropDAOImpl) GetAllCrops(ctx context.Context) ([]models.CropProfile, error) {
	query := `
		SELECT id, name, min_temp, max_temp, min_precipitation, max_precipitation, ideal_season 
		FROM crops
	`
	rows, err := dao.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var crops []models.CropProfile
	for rows.Next() {
		var c models.CropProfile
		if err := rows.Scan(&c.ID, &c.Name, &c.MinTemp, &c.MaxTemp, &c.MinPrecipitation, &c.MaxPrecipitation, &c.IdealSeason); err != nil {
			return nil, err
		}
		crops = append(crops, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return crops, nil
}
