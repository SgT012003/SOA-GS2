package main

import (
	"log"

	"agri-ai-api/internal/dao"
	"agri-ai-api/internal/handlers"
	"agri-ai-api/internal/middleware"

	_ "agri-ai-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Agri-AI API
// @version 1.0
// @description API backend for agricultural AI insights
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Initialize database
	if err := dao.InitDB(); err != nil {
		log.Printf("Warning: Database initialization failed (is docker running?): %v", err)
	} else {
		defer dao.CloseDB()
	}

	// Initialize Gin engine
	r := gin.Default()

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", handlers.PingHandler)

		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", handlers.RegisterHandler)
			authGroup.POST("/login", handlers.LoginHandler)
		}

		// Protected routes
		protected := v1.Group("/protected")
		protected.Use(middleware.AuthMiddleware())
		protected.Use(middleware.RateLimitMiddleware())
		protected.Use(middleware.UsageLogMiddleware())
		{
			protected.GET("/me", func(c *gin.Context) {
				userID := c.MustGet("userID").(int)
				c.JSON(200, gin.H{
					"message": "You have access!",
					"user_id": userID,
				})
			})
			protected.GET("/weather", handlers.GetWeatherHandler)
			
			engine := protected.Group("/engine")
			{
				engine.GET("/harvest", handlers.HarvestEngineHandler)
			}
		}
	}

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
