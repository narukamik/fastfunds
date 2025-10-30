package handlers

import (
	"fastfunds/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	accountService *service.AccountService,
	transactionService *service.TransactionService,
) {
	accountHandler := NewAccountHandler(accountService)
	transactionHandler := NewTransactionHandler(transactionService)

	router.POST("/accounts", accountHandler.CreateAccount)

	router.GET("/accounts/:account_id", accountHandler.GetAccount)
	router.POST("/transactions", transactionHandler.SubmitTransaction)
}
