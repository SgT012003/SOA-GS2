package models

// ProblemDetail implementa a especificação IETF RFC 7807 para padronização de erros em APIs HTTP.
// O campo 'business_error' foi adicionado como extensão para facilitar a identificação pelos frontends.
type ProblemDetail struct {
	Type          string `json:"type" example:"https://api.agri-ai.com/errors/invalid-parameters"`
	Title         string `json:"title" example:"Parâmetros Inválidos"`
	Status        int    `json:"status" example:"400"`
	Detail        string `json:"detail" example:"A latitude fornecida não é um número válido."`
	Instance      string `json:"instance,omitempty" example:"/api/v1/protected/engine/harvest"`
	BusinessError bool   `json:"business_error" example:"true"`
}

// NewBusinessError cria uma estrutura ProblemDetail focada em regras de negócio não atendidas.
func NewBusinessError(title, detail, instance string, status int) ProblemDetail {
	return ProblemDetail{
		Type:          "about:blank", // Ou URL com documentação do erro
		Title:         title,
		Status:        status,
		Detail:        detail,
		Instance:      instance,
		BusinessError: true,
	}
}

// NewTechnicalError cria uma estrutura ProblemDetail focada em falhas técnicas internas (banco, integração, parsing).
func NewTechnicalError(title, detail, instance string, status int) ProblemDetail {
	return ProblemDetail{
		Type:          "about:blank",
		Title:         title,
		Status:        status,
		Detail:        detail,
		Instance:      instance,
		BusinessError: false,
	}
}
