package handlers

import (
	"net/http"

	"agri-ai-api/internal/auth"
	"agri-ai-api/internal/dao"
	"agri-ai-api/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var userDAO = dao.NewUserDAO()

// RegisterHandler cria uma nova conta de usuário
// @Summary Registrar usuário
// @Description Cria uma nova conta com e-mail e senha
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Credenciais do usuário"
// @Success 201 {object} map[string]interface{} "Usuário criado com sucesso"
// @Failure 400 {object} map[string]string "Requisição inválida"
// @Failure 409 {object} map[string]string "Usuário já existe"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /auth/register [post]
func RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Verificar se já existe
	existingUser, err := userDAO.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := userDAO.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": newUser.ID,
	})
}

// LoginHandler autentica um usuário e retorna o JWT
// @Summary Login do usuário
// @Description Autentica e-mail/senha e retorna um JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Credenciais do usuário"
// @Success 200 {object} models.LoginResponse "Token JWT"
// @Failure 400 {object} map[string]string "Requisição inválida"
// @Failure 401 {object} map[string]string "Credenciais inválidas"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /auth/login [post]
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Buscar usuário
	user, err := userDAO.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Gerar JWT
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
	})
}
