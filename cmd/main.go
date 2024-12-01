// @title Financial Tracker API
// @version 1.0
// @description This is a simple API for a financial tracking system.

// @contact.name API Support
// @contact.email alexey_gladilin@mail.ru

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nemopss/financial-tracker/config"
	_ "github.com/nemopss/financial-tracker/docs" // Swagger docs
	"github.com/nemopss/financial-tracker/internal/handlers"
	"github.com/nemopss/financial-tracker/internal/middleware"
	"github.com/nemopss/financial-tracker/internal/repository"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	if cfg.JWTSecret == "" || cfg.Port == "" || cfg.DBHost == "" {
		log.Fatal("Missing critical configuration values")
	}

	// Connect to the database
	db, err := repository.NewDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Conn.Close()

	// Initialize handlers
	authHandler := &handlers.AuthHandler{Repo: db, JWTSecret: cfg.JWTSecret}
	categoryHandler := &handlers.CategoryHandler{Repo: db}
	transactionHandler := &handlers.TransactionHandler{Repo: db}
	analyticsHandler := &handlers.AnalyticsHandler{Repo: db}

	// Initialize Gin
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		api.POST("/auth/register", authHandler.RegisterGin)
		api.POST("/auth/login", authHandler.LoginGin)

		// Protected routes
		protected := api.Group("/", middleware.AuthGin(cfg.JWTSecret))
		{
			// Categories
			protected.POST("/categories", categoryHandler.CreateCategoryGin)
			protected.GET("/categories/list", categoryHandler.GetCategoriesGin)
			protected.PUT("/categories/update", categoryHandler.UpdateCategoryGin)
			protected.DELETE("/categories/delete", categoryHandler.DeleteCategoryGin)

			// Transactions
			protected.POST("/transactions", transactionHandler.CreateTransactionGin)
			protected.GET("/transactions/list", transactionHandler.GetTransactionsGin)
			protected.PUT("/transactions/update", transactionHandler.UpdateTransactionGin)
			protected.DELETE("/transactions/delete", transactionHandler.DeleteTransactionGin)

			// Analytics
			protected.GET("/analytics/income-expenses", analyticsHandler.GetIncomeAndExpensesGin)
			protected.GET("/analytics/categories", analyticsHandler.GetCategoryAnalyticsGin)
			protected.GET("/analytics/income-expenses-filtered", analyticsHandler.GetIncomeAndExpensesFilteredGin)
			protected.GET("/analytics/categories-filtered", analyticsHandler.GetCategoryAnalyticsFilteredGin)
		}
	}

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
