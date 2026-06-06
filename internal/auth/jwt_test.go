package auth

import (
	"testing"
)

func TestGenerateAndValidateToken(t *testing.T) {
	expectedUserID := 123

	tokenStr, err := GenerateToken(expectedUserID)
	if err != nil {
		t.Fatalf("falha ao gerar token: %v", err)
	}

	if tokenStr == "" {
		t.Fatal("token gerado está vazio")
	}

	claims, err := ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("falha ao validar token: %v", err)
	}

	if claims.UserID != expectedUserID {
		t.Errorf("esperado userID %d, recebido %d", expectedUserID, claims.UserID)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	_, err := ValidateToken("invalid.token.string")
	if err == nil {
		t.Fatal("esperado erro ao validar token inválido, recebido nil")
	}
}
