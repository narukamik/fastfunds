package service

import (
	"database/sql"
	"errors"
	"fastfunds/internal/models"
	"fastfunds/internal/repository"
	"fastfunds/internal/util"
	"time"
)

func NewTransactionService(
	db *sql.DB,
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
) *TransactionService {
	return &TransactionService{
		db:              db,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

type TransactionService struct {
	db              *sql.DB
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func (s *TransactionService) ProcessTransaction(req *models.TransactionRequest) error {
	// Validate request
	if req.SourceAccountID <= 0 || req.DestinationAccountID <= 0 {
		return errors.New("invalid account IDs")
	}

	if req.SourceAccountID == req.DestinationAccountID {
		return errors.New("source and destination accounts cannot be the same")
	}

	if req.Amount == "" {
		return errors.New("amount is required")
	}

	// Validate and convert amount to pennies
	amountPennies, err := util.DecimalStringToPennies(req.Amount)
	if err != nil || amountPennies <= 0 {
		return errors.New("invalid amount format")
	}

	// Start DB transaction
	tx, err := s.db.Begin()

	if err != nil {
		return errors.New("couldn't start DB transaction")
	}

	defer tx.Rollback()

	// Get source account
	sourceAccount, err := s.accountRepo.SelectTx(tx, req.SourceAccountID)
	if err != nil {
		return errors.New("source account not found")
	}

	// Get destination account
	destAccount, err := s.accountRepo.SelectTx(tx, req.DestinationAccountID)
	if err != nil {
		return errors.New("destination account not found")
	}

	// Check source account balance
	if sourceAccount.CurrentBalance < amountPennies {
		return errors.New("insufficient funds")
	}

	// Calculate new balances in pennies
	newSourceBalance := sourceAccount.CurrentBalance - amountPennies
	newDestBalance := destAccount.CurrentBalance + amountPennies

	// Update accounts
	sourceAccount.CurrentBalance = newSourceBalance
	destAccount.CurrentBalance = newDestBalance

	if err := s.accountRepo.UpdateTx(tx, sourceAccount); err != nil {
		return errors.New("failed to update source account")
	}

	if err := s.accountRepo.UpdateTx(tx, destAccount); err != nil {
		return errors.New("failed to update destination account")
	}

	// Create transaction record
	transaction := &models.Transaction{
		SourceAccountID:      req.SourceAccountID,
		DestinationAccountID: req.DestinationAccountID,
		AmountPennies:        amountPennies,
		Status:               "completed",
		CreatedAt:            time.Now().Format(time.RFC3339),
	}

	if err := s.transactionRepo.CreateTx(tx, transaction); err != nil {
		return errors.New("transaction creation failed")
	}

	if err = tx.Commit(); err != nil {
		return errors.New("couldn't commit db transaction")
	}

	return nil
}
