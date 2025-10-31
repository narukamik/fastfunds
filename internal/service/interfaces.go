package service

import "fastfunds/internal/models"

type IAccountService interface {
	CreateAccount(req *models.CreateAccountRequest) error
	GetAccount(accountID int) (*models.AccountView, error)
}

type ITransactionService interface {
	ProcessTransaction(req *models.TransactionRequest) error
}
