package tests

import (
	"database/sql"
	"errors"
	"fastfunds/internal/models"
	"fastfunds/internal/service"
	"fastfunds/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockAccountRepository struct {
	createFn   func(*models.Account) error
	getByIDFn  func(int) (*models.Account, error)
	existsFn   func(int) (bool, error)
	selectTxFn func(*sql.Tx, int) (*models.Account, error)
	updateTxFn func(*sql.Tx, *models.Account) error
}

func (m *mockAccountRepository) Create(a *models.Account) error {
	if m.createFn != nil {
		return m.createFn(a)
	}
	return nil
}

func (m *mockAccountRepository) GetByID(id int) (*models.Account, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(id)
	}
	return nil, nil
}

func (m *mockAccountRepository) Exists(id int) (bool, error) {
	if m.existsFn != nil {
		return m.existsFn(id)
	}
	return false, nil
}

func (m *mockAccountRepository) SelectTx(tx *sql.Tx, id int) (*models.Account, error) {
	if m.selectTxFn != nil {
		return m.selectTxFn(tx, id)
	}
	return nil, nil
}

func (m *mockAccountRepository) UpdateTx(tx *sql.Tx, a *models.Account) error {
	if m.updateTxFn != nil {
		return m.updateTxFn(tx, a)
	}
	return nil
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

func TestCreateAccount(t *testing.T) {
	cases := []struct {
		name    string
		req     *models.CreateAccountRequest
		money   util.MoneyConverter
		repo    *mockAccountRepository
		wantErr string
	}{
		{
			name:    "invalid_id",
			req:     &models.CreateAccountRequest{AccountID: 0, InitialBalance: "10.00"},
			money:   &mockMoneyConverter{},
			repo:    &mockAccountRepository{},
			wantErr: "invalid account_id",
		},
		{
			name:    "empty_balance",
			req:     &models.CreateAccountRequest{AccountID: 1, InitialBalance: ""},
			money:   &mockMoneyConverter{},
			repo:    &mockAccountRepository{},
			wantErr: "initial balance is required",
		},
		{
			name: "invalid_balance_format",
			req:  &models.CreateAccountRequest{AccountID: 1, InitialBalance: "abc"},
			money: &mockMoneyConverter{
				decFn: func(s string) (int64, error) { return 0, errors.New("bad") },
			},
			repo:    &mockAccountRepository{},
			wantErr: "invalid balance format",
		},
		{
			name: "exists_error",
			req:  &models.CreateAccountRequest{AccountID: 2, InitialBalance: "1.00"},
			money: &mockMoneyConverter{
				decFn: func(s string) (int64, error) { return 100, nil },
			},
			repo: &mockAccountRepository{
				existsFn: func(int) (bool, error) { return false, errors.New("db") },
			},
			wantErr: "db",
		},
		{
			name: "already_exists",
			req:  &models.CreateAccountRequest{AccountID: 3, InitialBalance: "1.00"},
			money: &mockMoneyConverter{
				decFn: func(s string) (int64, error) { return 100, nil },
			},
			repo: &mockAccountRepository{
				existsFn: func(int) (bool, error) { return true, nil },
			},
			wantErr: "account already exists",
		},
		{
			name: "create_error",
			req:  &models.CreateAccountRequest{AccountID: 4, InitialBalance: "1.00"},
			money: &mockMoneyConverter{
				decFn: func(s string) (int64, error) { return 100, nil },
			},
			repo: &mockAccountRepository{
				existsFn: func(int) (bool, error) { return false, nil },
				createFn: func(*models.Account) error { return errors.New("fail") },
			},
			wantErr: "fail",
		},
		{
			name: "success",
			req:  &models.CreateAccountRequest{AccountID: 5, InitialBalance: "123.45"},
			money: &mockMoneyConverter{
				decFn: func(s string) (int64, error) { return 12345, nil },
			},
			repo: &mockAccountRepository{
				existsFn: func(int) (bool, error) { return false, nil },
				createFn: func(a *models.Account) error {
					assert.Equal(t, 5, a.AccountID)
					assert.Equal(t, int64(12345), a.CurrentBalance)
					return nil
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewAccountServiceWithDeps(tc.repo, tc.money)
			err := svc.CreateAccount(tc.req)
			if tc.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetAccount(t *testing.T) {
	cases := []struct {
		name     string
		id       int
		repo     *mockAccountRepository
		money    util.MoneyConverter
		wantErr  string
		wantView *models.AccountView
	}{
		{
			name:    "invalid_id",
			id:      0,
			repo:    &mockAccountRepository{},
			money:   &mockMoneyConverter{},
			wantErr: "invalid account_id",
		},
		{
			name: "repo_error",
			id:   10,
			repo: &mockAccountRepository{
				getByIDFn: func(int) (*models.Account, error) { return nil, errors.New("x") },
			},
			money:   &mockMoneyConverter{},
			wantErr: "couldn't get account by ID",
		},
		{
			name: "success_negative_format",
			id:   33,
			repo: &mockAccountRepository{
				getByIDFn: func(id int) (*models.Account, error) {
					return &models.Account{AccountID: id, CurrentBalance: -123}, nil
				},
			},
			money: &mockMoneyConverter{
				fmtFn: func(p int64) string {
					assert.Equal(t, int64(-123), p)
					return "-1.23"
				},
			},
			wantView: &models.AccountView{AccountID: 33, CurrentBalance: "-1.23"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewAccountServiceWithDeps(tc.repo, tc.money)
			got, err := svc.GetAccount(tc.id)
			if tc.wantErr != "" {
				assert.Nil(t, got)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantView, got)
		})
	}
}
