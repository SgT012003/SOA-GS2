package services

import (
	"errors"

	"agri-ai-api/internal/auth"
	"agri-ai-api/internal/dao"
	"agri-ai-api/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// AuthService define as operações de negócio para autenticação
type AuthService interface {
	RegisterUser(req models.RegisterRequest) (*models.User, error)
	LoginUser(req models.LoginRequest) (string, error)
}

// AuthServiceImpl implementa AuthService
type AuthServiceImpl struct {
	userDAO dao.UserDAO
}

// NewAuthService cria uma nova instância
func NewAuthService(userDAO dao.UserDAO) AuthService {
	return &AuthServiceImpl{
		userDAO: userDAO,
	}
}

// RegisterUser lida com o registro de usuários
func (s *AuthServiceImpl) RegisterUser(req models.RegisterRequest) (*models.User, error) {
	existingUser, err := s.userDAO.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("failed to check existing user")
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	newUser := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userDAO.CreateUser(newUser); err != nil {
		return nil, errors.New("failed to create user")
	}

	return newUser, nil
}

// LoginUser lida com o login e retorna um JWT
func (s *AuthServiceImpl) LoginUser(req models.LoginRequest) (string, error) {
	user, err := s.userDAO.GetUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("failed to fetch user")
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
