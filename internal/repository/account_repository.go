package repository

import (
	"database/sql"
	"errors"
	"fastfunds/internal/models"
)

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

type PostgresAccountRepository struct {
	db *sql.DB
}

func (r *PostgresAccountRepository) Create(account *models.Account) error {
	return r.db.QueryRow(
		`INSERT INTO accounts (account_id, balance) VALUES ($1, $2) RETURNING account_id`,
		account.AccountID, account.CurrentBalance,
	).Scan(&account.AccountID)
}

// lock to void lost update issue
func (r *PostgresAccountRepository) SelectTx(tx *sql.Tx, id int) (*models.Account, error) {
	acc := &models.Account{}
	row := tx.QueryRow(
		`SELECT account_id, balance FROM accounts WHERE account_id = $1 FOR UPDATE`, id,
	)
	if err := row.Scan(&acc.AccountID, &acc.CurrentBalance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return acc, nil
}

func (r *PostgresAccountRepository) GetByID(id int) (*models.Account, error) {
	acc := &models.Account{}
	row := r.db.QueryRow(
		`SELECT account_id, balance FROM accounts WHERE account_id = $1`, id,
	)
	if err := row.Scan(&acc.AccountID, &acc.CurrentBalance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		} else {
			return nil, err
		}
	}
	return acc, nil
}

func (r *PostgresAccountRepository) UpdateTx(tx *sql.Tx, account *models.Account) error {
	res, err := tx.Exec(
		`UPDATE accounts SET balance = $2 WHERE account_id = $1`,
		account.AccountID, account.CurrentBalance,
	)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("account not found")
	}
	return nil
}

func (r *PostgresAccountRepository) Exists(id int) (bool, error) {
	var exists bool
	if err := r.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM accounts WHERE account_id = $1)`, id).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
