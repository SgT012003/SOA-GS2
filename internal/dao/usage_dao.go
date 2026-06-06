package dao

import (
	"database/sql"
)

// UsageDAO define os métodos de acesso aos dados de logs de uso
type UsageDAO interface {
	LogUsage(userID int, endpoint string) error
}

// UsageDAOImpl implementação do UsageDAO
type UsageDAOImpl struct {
	db *sql.DB
}

// NewUsageDAO retorna a implementação
func NewUsageDAO() *UsageDAOImpl {
	return &UsageDAOImpl{db: DB}
}

// LogUsage insere um registro de uso da API
func (dao *UsageDAOImpl) LogUsage(userID int, endpoint string) error {
	query := `INSERT INTO api_usage_logs (user_id, endpoint) VALUES ($1, $2)`
	_, err := dao.db.Exec(query, userID, endpoint)
	return err
}
