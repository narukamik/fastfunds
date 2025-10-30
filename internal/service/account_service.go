package service

import (
	"errors"
	"fastfunds/internal/models"
	"fastfunds/internal/repository"
	"fastfunds/internal/util"
)

func NewAccountService(accountRepo repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

type AccountService struct {
	accountRepo repository.AccountRepository
}

func (s *AccountService) CreateAccount(req *models.CreateAccountRequest) error {
	if req.AccountID <= 0 {
		return errors.New("invalid account_id")
	}

	if req.InitialBalance == "" {
		return errors.New("initial balance is required")
	}

	pennies, err := util.DecimalStringToPennies(req.InitialBalance)
	if err != nil {
		return errors.New("invalid balance format")
	}

	exists, err := s.accountRepo.Exists(req.AccountID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("account already exists")
	}

	account := &models.Account{
		AccountID:      req.AccountID,
		CurrentBalance: pennies,
	}

	return s.accountRepo.Create(account)
}

func (s *AccountService) GetAccount(accountID int) (*models.AccountView, error) {
	if accountID <= 0 {
		return nil, errors.New("invalid account_id")
	}

	account, err := s.accountRepo.GetByID(accountID)

	if err != nil {
		return nil, errors.New("couldn't get account by ID")
	}

	accountView := &models.AccountView{
		AccountID:      account.AccountID,
		CurrentBalance: util.PenniesToDecimalString(account.CurrentBalance),
	}

	return accountView, nil
}
