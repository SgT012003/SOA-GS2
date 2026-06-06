package dao

import (
	"database/sql"
	"errors"

	"agri-ai-api/internal/models"
)

// UserDAO define os métodos de acesso aos dados de usuário
type UserDAO interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
}

// UserDAOImpl é a implementação do UserDAO
type UserDAOImpl struct {
	db *sql.DB
}

// NewUserDAO retorna uma nova instância de UserDAOImpl
func NewUserDAO() *UserDAOImpl {
	return &UserDAOImpl{db: DB}
}

// CreateUser insere um novo usuário no banco de dados
func (dao *UserDAOImpl) CreateUser(user *models.User) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := dao.db.QueryRow(query, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	return err
}

// GetUserByEmail busca um usuário pelo email
func (dao *UserDAOImpl) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	err := dao.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return user, nil
}

// GetUserByID busca um usuário pelo ID
func (dao *UserDAOImpl) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, created_at FROM users WHERE id = $1`
	err := dao.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return user, nil
}
