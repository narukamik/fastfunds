package repository

import (
	"database/sql"
	"errors"
	"fastfunds/internal/models"
)

func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

type PostgresTransactionRepository struct {
	db *sql.DB
}

func (r *PostgresTransactionRepository) CreateTx(tx *sql.Tx, t *models.Transaction) error {
	return tx.QueryRow(
		`INSERT INTO transactions (source_account_id, destination_account_id, amount, status)
         VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		t.SourceAccountID, t.DestinationAccountID, t.AmountPennies, t.Status,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *PostgresTransactionRepository) GetByID(id int) (*models.Transaction, error) {
	t := &models.Transaction{}
	row := r.db.QueryRow(
		`SELECT id, source_account_id, destination_account_id, amount, status, created_at
		 FROM transactions WHERE id = $1`, id,
	)
	if err := row.Scan(&t.ID, &t.SourceAccountID, &t.DestinationAccountID, &t.AmountPennies, &t.Status, &t.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return t, nil
}

func (r *PostgresTransactionRepository) GetByAccountID(accountID int) ([]*models.Transaction, error) {
	rows, err := r.db.Query(
		`SELECT id, source_account_id, destination_account_id, amount, status, created_at
		 FROM transactions
		 WHERE source_account_id = $1 OR destination_account_id = $1
		 ORDER BY id DESC`, accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Transaction
	for rows.Next() {
		t := &models.Transaction{}
		if err := rows.Scan(&t.ID, &t.SourceAccountID, &t.DestinationAccountID, &t.AmountPennies, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}
