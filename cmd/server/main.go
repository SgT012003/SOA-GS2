package main

import (
	"log/slog"
	"os"

	"agri-ai-api/internal/dao"
	"agri-ai-api/internal/handlers"
	"agri-ai-api/internal/middleware"
	"agri-ai-api/internal/services"

	_ "agri-ai-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
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
	// Configurar slog (Structured Logging em JSON)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize database
	if err := dao.InitDB(); err != nil {
		slog.Warn("Database initialization failed (is docker running?)", slog.String("error", err.Error()))
	} else {
		defer dao.CloseDB()
	}

	// Wiring (Injeção de Dependências)
	userDAO := dao.NewUserDAO()
	weatherDAO := dao.NewWeatherDAO()
	usageDAO := dao.NewUsageDAO()
	cropDAO := dao.NewCropDAO(dao.DB) // Injetando dao.DB global

	authService := services.NewAuthService(userDAO)
	weatherService := services.NewWeatherService(weatherDAO)
	engineService := services.NewEngineService(weatherService, cropDAO)
	cropService := services.NewCropService(cropDAO)

	authHandler := handlers.NewAuthHandler(authService)
	weatherHandler := handlers.NewWeatherHandler(weatherService)
	engineHandler := handlers.NewEngineHandler(engineService)
	cropHandler := handlers.NewCropHandler(cropService)

	// Initialize Gin engine
	r := gin.Default()

	// Configurar Prometheus Metrics
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/healthz", handlers.HealthzHandler)

		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := v1.Group("/protected")
		protected.Use(middleware.AuthMiddleware())
		protected.Use(middleware.RateLimitMiddleware())
		protected.Use(middleware.UsageLogMiddleware(usageDAO))
		{
			protected.GET("/me", func(c *gin.Context) {
				userID := c.MustGet("userID").(int)
				c.JSON(200, gin.H{
					"message": "You have access!",
					"user_id": userID,
				})
			})

			// Rotas de Clima
			protected.GET("/weather", weatherHandler.GetWeatherHandler)
			protected.GET("/weather/cache", weatherHandler.GetAllCachesHandler)

			// Rotas de Culturas
			protected.GET("/crops", cropHandler.GetAllCropsHandler)

			// Rotas do Motor Preditivo
			engine := protected.Group("/engine")
			{
				engine.GET("/harvest", engineHandler.HarvestEngineHandler)
				engine.GET("/risk-analysis", engineHandler.RiskAnalysisHandler)
				engine.GET("/crop-selector", engineHandler.CropSelectorHandler)
			}
		}
	}

	slog.Info("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		slog.Error("Server failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
