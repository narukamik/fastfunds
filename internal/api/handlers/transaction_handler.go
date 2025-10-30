package handlers

import (
	"fastfunds/internal/models"
	"fastfunds/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// SubmitTransaction godoc
// @Summary Submit transaction
// @Accept json
// @Produce json
// @Param request body models.TransactionRequest true "Transaction payload"
// @Success 200
// @Failure 400 {object} map[string]string
// @Router /transactions [post]
// @Tags transactions
func (h *TransactionHandler) SubmitTransaction(c *gin.Context) {
	var req models.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := h.transactionService.ProcessTransaction(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated) // 201
}
