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
	s := &TransactionService{
		db:              db,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		money:           util.DefaultMoneyConverter{},
	}
	s.beginFn = func() (*sql.Tx, error) { return s.db.Begin() }
	s.rollbackFn = func(tx *sql.Tx) error { return tx.Rollback() }
	s.commitFn = func(tx *sql.Tx) error { return tx.Commit() }
	return s
}

// NewTransactionServiceWithDeps allows injecting MoneyConverter for testing.
// NewTransactionServiceWithDeps allows injecting MoneyConverter and transaction functions for testing.
func NewTransactionServiceWithDeps(
	db *sql.DB,
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	money util.MoneyConverter,
	opts ...func(*TransactionService),
) *TransactionService {
	if money == nil {
		money = util.DefaultMoneyConverter{}
	}
	s := &TransactionService{
		db:              db,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		money:           money,
	}
	s.beginFn = func() (*sql.Tx, error) { return s.db.Begin() }
	s.rollbackFn = func(tx *sql.Tx) error { return tx.Rollback() }
	s.commitFn = func(tx *sql.Tx) error { return tx.Commit() }
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type TransactionService struct {
	db              *sql.DB
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	money           util.MoneyConverter
	beginFn         func() (*sql.Tx, error)
	rollbackFn      func(*sql.Tx) error
	commitFn        func(*sql.Tx) error
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
	amountPennies, err := s.money.DecimalStringToPennies(req.Amount)
	if err != nil || amountPennies <= 0 {
		return errors.New("invalid amount format")
	}

	// Start DB transaction
	tx, err := s.beginFn()

	if err != nil {
		return errors.New("couldn't start DB transaction")
	}
	if tx == nil {
		return errors.New("couldn't start DB transaction")
	}

	defer s.rollbackFn(tx)

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

	if err = s.commitFn(tx); err != nil {
		return errors.New("couldn't commit db transaction")
	}

	return nil
}
