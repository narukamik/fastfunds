package models

type Transaction struct {
	ID                   int    `json:"id"`
	SourceAccountID      int    `json:"source_account_id"`
	DestinationAccountID int    `json:"destination_account_id"`
	AmountPennies        int64  `json:"amount_pennies"`
	Status               string `json:"status"`
	CreatedAt            string `json:"created_at"`
}

type TransactionRequest struct {
	SourceAccountID      int    `json:"source_account_id"`
	DestinationAccountID int    `json:"destination_account_id"`
	Amount               string `json:"amount"`
}
