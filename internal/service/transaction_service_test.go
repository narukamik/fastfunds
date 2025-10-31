package service

import (
	"database/sql"
	"errors"
	"fastfunds/internal/models"
	"testing"
)

// SetBeginFn allows tests to override the beginFn for TransactionService
func (s *TransactionService) SetBeginFn(fn func() (*sql.Tx, error)) {
	s.beginFn = fn
}

// SetRollbackFn allows tests to override the rollbackFn for TransactionService
func (s *TransactionService) SetRollbackFn(fn func(*sql.Tx) error) {
	s.rollbackFn = fn
}

// SetCommitFn allows tests to override the commitFn for TransactionService
func (s *TransactionService) SetCommitFn(fn func(*sql.Tx) error) {
	s.commitFn = fn
}

type mockAccountRepo struct {
	SelectTxFunc func(tx *sql.Tx, id int) (*models.Account, error)
	UpdateTxFunc func(tx *sql.Tx, account *models.Account) error
}

func (m *mockAccountRepo) Create(account *models.Account) error    { return nil }
func (m *mockAccountRepo) GetByID(id int) (*models.Account, error) { return nil, nil }
func (m *mockAccountRepo) SelectTx(tx *sql.Tx, id int) (*models.Account, error) {
	if m.SelectTxFunc != nil {
		return m.SelectTxFunc(tx, id)
	}
	return nil, nil
}
func (m *mockAccountRepo) UpdateTx(tx *sql.Tx, account *models.Account) error {
	if m.UpdateTxFunc != nil {
		return m.UpdateTxFunc(tx, account)
	}
	return nil
}
func (m *mockAccountRepo) Exists(id int) (bool, error) { return true, nil }

type mockTransactionRepo struct {
	CreateTxFunc func(tx *sql.Tx, transaction *models.Transaction) error
}

func (m *mockTransactionRepo) CreateTx(tx *sql.Tx, transaction *models.Transaction) error {
	if m.CreateTxFunc != nil {
		return m.CreateTxFunc(tx, transaction)
	}
	return nil
}
func (m *mockTransactionRepo) GetByID(id int) (*models.Transaction, error) { return nil, nil }
func (m *mockTransactionRepo) GetByAccountID(accountID int) ([]*models.Transaction, error) {
	return nil, nil
}

type mockMoneyConverter struct {
	decFn func(string) (int64, error)
	fmtFn func(int64) string
}

func (m mockMoneyConverter) DecimalStringToPennies(s string) (int64, error) {
	if m.decFn != nil {
		return m.decFn(s)
	}
	return 0, nil
}

func (m mockMoneyConverter) PenniesToDecimalString(p int64) string {
	if m.fmtFn != nil {
		return m.fmtFn(p)
	}
	return ""
}

func setTxnFns(ts *TransactionService) {
	ts.SetBeginFn(func() (*sql.Tx, error) { return &sql.Tx{}, nil })
	ts.SetRollbackFn(func(tx *sql.Tx) error { return nil })
	ts.SetCommitFn(func(tx *sql.Tx) error { return nil })
}

func TestProcessTransaction_Success(t *testing.T) {
	accountRepo := &mockAccountRepo{
		SelectTxFunc: func(tx *sql.Tx, id int) (*models.Account, error) {
			if id == 1 {
				return &models.Account{AccountID: 1, CurrentBalance: 1000}, nil
			}
			if id == 2 {
				return &models.Account{AccountID: 2, CurrentBalance: 500}, nil
			}
			return nil, errors.New("not found")
		},
		UpdateTxFunc: func(tx *sql.Tx, account *models.Account) error { return nil },
	}
	transactionRepo := &mockTransactionRepo{
		CreateTxFunc: func(tx *sql.Tx, transaction *models.Transaction) error { return nil },
	}
	money := &mockMoneyConverter{
		decFn: func(s string) (int64, error) { return 200, nil },
		fmtFn: func(p int64) string { return "2.00" },
	}

	ts := NewTransactionServiceWithDeps(&sql.DB{}, accountRepo, transactionRepo, money)
	setTxnFns(ts)

	req := &models.TransactionRequest{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               "2.00",
	}

	err := ts.ProcessTransaction(req)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestProcessTransaction_InvalidAmount(t *testing.T) {
	money := &mockMoneyConverter{
		decFn: func(s string) (int64, error) { return 0, errors.New("bad format") },
	}
	ts := NewTransactionServiceWithDeps(&sql.DB{}, &mockAccountRepo{}, &mockTransactionRepo{}, money)
	setTxnFns(ts)

	req := &models.TransactionRequest{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               "bad",
	}

	err := ts.ProcessTransaction(req)
	if err == nil || err.Error() != "invalid amount format" {
		t.Errorf("expected invalid amount format error, got: %v", err)
	}
}

func TestProcessTransaction_InsufficientFunds(t *testing.T) {
	accountRepo := &mockAccountRepo{
		SelectTxFunc: func(tx *sql.Tx, id int) (*models.Account, error) {
			if id == 1 {
				return &models.Account{AccountID: 1, CurrentBalance: 100}, nil
			}
			if id == 2 {
				return &models.Account{AccountID: 2, CurrentBalance: 500}, nil
			}
			return nil, errors.New("not found")
		},
		UpdateTxFunc: func(tx *sql.Tx, account *models.Account) error { return nil },
	}
	money := &mockMoneyConverter{
		decFn: func(s string) (int64, error) { return 200, nil },
	}
	ts := NewTransactionServiceWithDeps(&sql.DB{}, accountRepo, &mockTransactionRepo{}, money)
	setTxnFns(ts)

	req := &models.TransactionRequest{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               "2.00",
	}

	err := ts.ProcessTransaction(req)
	if err == nil || err.Error() != "insufficient funds" {
		t.Errorf("expected insufficient funds error, got: %v", err)
	}
}

func TestProcessTransaction_SameAccount(t *testing.T) {
	ts := NewTransactionServiceWithDeps(&sql.DB{}, &mockAccountRepo{}, &mockTransactionRepo{}, &mockMoneyConverter{})
	setTxnFns(ts)

	req := &models.TransactionRequest{
		SourceAccountID:      1,
		DestinationAccountID: 1,
		Amount:               "2.00",
	}

	err := ts.ProcessTransaction(req)
	if err == nil || err.Error() != "source and destination accounts cannot be the same" {
		t.Errorf("expected same account error, got: %v", err)
	}
}

func TestProcessTransaction_SourceAccountNotFound(t *testing.T) {
	accountRepo := &mockAccountRepo{
		SelectTxFunc: func(tx *sql.Tx, id int) (*models.Account, error) {
			if id == 1 {
				return nil, errors.New("not found")
			}
			return &models.Account{AccountID: 2, CurrentBalance: 500}, nil
		},
		UpdateTxFunc: func(tx *sql.Tx, account *models.Account) error { return nil },
	}
	ts := NewTransactionServiceWithDeps(&sql.DB{}, accountRepo, &mockTransactionRepo{}, &mockMoneyConverter{decFn: func(s string) (int64, error) { return 200, nil }})
	setTxnFns(ts)

	req := &models.TransactionRequest{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               "2.00",
	}

	err := ts.ProcessTransaction(req)
	if err == nil || err.Error() != "source account not found" {
		t.Errorf("expected source account not found error, got: %v", err)
	}
}

func TestProcessTransaction_DestinationAccountNotFound(t *testing.T) {
	accountRepo := &mockAccountRepo{
		SelectTxFunc: func(tx *sql.Tx, id int) (*models.Account, error) {
			if id == 1 {
				return &models.Account{AccountID: 1, CurrentBalance: 1000}, nil
			}
			return nil, errors.New("not found")
		},
		UpdateTxFunc: func(tx *sql.Tx, account *models.Account) error { return nil },
	}
	ts := NewTransactionServiceWithDeps(&sql.DB{}, accountRepo, &mockTransactionRepo{}, &mockMoneyConverter{decFn: func(s string) (int64, error) { return 200, nil }})
	setTxnFns(ts)

	req := &models.TransactionRequest{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               "2.00",
	}

	err := ts.ProcessTransaction(req)
	if err == nil || err.Error() != "destination account not found" {
		t.Errorf("expected destination account not found error, got: %v", err)
	}
}
