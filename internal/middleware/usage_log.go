package middleware

import (
	"log/slog"

	"agri-ai-api/internal/dao"

	"github.com/gin-gonic/gin"
)

// UsageLogMiddleware registra endpoints acessados por usuários autenticados
func UsageLogMiddleware(usageDAO dao.UsageDAO) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Executa os handlers primeiro
		c.Next()

		// Apenas registrar requisições com sucesso
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			userID, exists := c.Get("userID")
			if exists {
				endpoint := c.Request.URL.Path

				// Log de forma assíncrona para não atrasar a resposta ao cliente
				go func(id int, ep string) {
					err := usageDAO.LogUsage(id, ep)
					if err != nil {
						slog.Error("Falha ao registrar log de uso", slog.String("error", err.Error()), slog.Int("user_id", id))
					}
				}(userID.(int), endpoint)
			}
		}
	}
}
