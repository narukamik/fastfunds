package handlers

import (
	"fastfunds/internal/models"
	"fastfunds/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewAccountHandler(accountService service.IAccountService) *AccountHandler {
	return &AccountHandler{
		 accountService: accountService,
	}
}

type AccountHandler struct {
	accountService service.IAccountService
}

// CreateAccount godoc
// @Summary Create account
// @Accept json
// @Produce json
// @Param request body models.CreateAccountRequest true "Create account payload"
// @Success 201
// @Failure 400 {object} map[string]string
// @Router /accounts [post]
// @Tags accounts
// Process account creation
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req models.CreateAccountRequest

	// Parse JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	if err := h.accountService.CreateAccount(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return empty response on success
	c.Status(http.StatusCreated)
}

// GetAccount godoc
// @Summary Get account information by ID
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {object} models.AccountView
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /accounts/{account_id} [get]
// @Tags accounts
// Get account
func (h *AccountHandler) GetAccount(c *gin.Context) {
	accountIDStr := c.Param("account_id")

	// Parse account ID from URL parameter
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account_id format"})
		return
	}

	account, err := h.accountService.GetAccount(accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}
