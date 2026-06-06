package middleware

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitMiddleware limita o número de chamadas por usuário autenticado
func RateLimitMiddleware() gin.HandlerFunc {
	rate, _ := limiter.NewRateFromFormatted("10-M") // 10 requisições por minuto
	store := memory.NewStore()
	limiterInstance := limiter.New(store, rate)

	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.Next()
			return
		}

		key := strconv.Itoa(userID.(int))
		context, err := limiterInstance.Get(c, key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			return
		}

		if context.Reached {
			slog.Warn("Rate limit excedido", slog.Int("user_id", userID.(int)), slog.String("path", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Muitas requisições. Aguarde um momento antes de tentar novamente.",
			})
			return
		}

		c.Next()
	}
}
