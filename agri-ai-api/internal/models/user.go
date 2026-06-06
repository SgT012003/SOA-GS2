package models

import "time"

// User representa a tabela users no banco de dados
type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Não retornar a senha em JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// RegisterRequest é o payload para criação de conta
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"strongpassword"`
}

// LoginRequest é o payload para autenticação
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"strongpassword"`
}

// LoginResponse é o payload de retorno após login bem sucedido
type LoginResponse struct {
	Token string `json:"token"`
}
