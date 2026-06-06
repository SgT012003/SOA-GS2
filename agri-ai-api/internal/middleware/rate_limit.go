package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var rateLimiter *limiter.Limiter

func init() {
	// Definir a taxa: 10 requisições por minuto
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  10,
	}

	// Usar um armazenamento em memória
	store := memory.NewStore()

	rateLimiter = limiter.New(store, rate)
}

// RateLimitMiddleware limita o número de chamadas por usuário autenticado
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.Next()
			return
		}

		key := strconv.Itoa(userID.(int))
		context, err := rateLimiter.Get(c, key)
		if err != nil {
			log.Printf("Erro no rate limiter: %v", err)
			c.Next() // Deixar passar em caso de falha no limitador para não travar a API
			return
		}

		if context.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Muitas requisições. Aguarde um momento antes de tentar novamente.",
			})
			return
		}

		c.Next()
	}
}
