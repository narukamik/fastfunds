package repository

import (
	"database/sql"
	"fastfunds/internal/models"
)

type AccountRepository interface {
	Create(account *models.Account) error
	GetByID(id int) (*models.Account, error)
	SelectTx(tx *sql.Tx, id int) (*models.Account, error)
	UpdateTx(tx *sql.Tx, account *models.Account) error
	Exists(id int) (bool, error)
}

type TransactionRepository interface {
	CreateTx(tx *sql.Tx, transaction *models.Transaction) error
	GetByID(id int) (*models.Transaction, error)
	GetByAccountID(accountID int) ([]*models.Transaction, error)
}
