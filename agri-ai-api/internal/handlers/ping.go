package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler responds with pong
// @Summary Ping the API
// @Description get a simple pong message to verify the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
