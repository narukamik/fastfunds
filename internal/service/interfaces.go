package service

import "fastfunds/internal/models"

// IAccountService defines the public contract for account operations.
type IAccountService interface {
    CreateAccount(req *models.CreateAccountRequest) error
    GetAccount(accountID int) (*models.AccountView, error)
}

// ITransactionService defines the public contract for transaction operations.
// Aligned to the current implementation in transaction_service.go.
type ITransactionService interface {
    ProcessTransaction(req *models.TransactionRequest) error
}
