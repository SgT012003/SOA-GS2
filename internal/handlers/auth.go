package handlers

import (
	"net/http"

	"agri-ai-api/internal/models"
	"agri-ai-api/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler encapsula os endpoints de autenticação
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler cria um novo AuthHandler
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register cria uma nova conta de usuário
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
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	newUser, err := h.authService.RegisterUser(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": newUser.ID,
	})
}

// Login autentica um usuário e retorna o JWT
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
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	token, err := h.authService.LoginUser(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "invalid email or password" {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
	})
}
