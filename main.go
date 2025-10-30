package main

import (
	"database/sql"
	_ "fastfunds/docs"
	"fastfunds/internal/api/handlers"
	"fastfunds/internal/repository"
	"fastfunds/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title FastFunds API
// @version 1.0
// @description FastFunds API server
// @host localhost:8080
// @BasePath /
func main() {
	dsn := os.Getenv("DATABASE_URL")
	log.Print("Database URL:", dsn)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Repositories init
	accountRepo := repository.NewPostgresAccountRepository(db)
	transactionRepo := repository.NewPostgresTransactionRepository(db)

	// Services init
	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(db, accountRepo, transactionRepo)

	// Init Gin router
	router := gin.Default()

	// Setup routes
	handlers.SetupRoutes(router, accountService, transactionService)

	// Setup Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
