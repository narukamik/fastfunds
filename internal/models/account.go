package models

type Account struct {
	AccountID      int   `json:"account_id"`
	CurrentBalance int64 `json:"current_balance"`
}

type AccountView struct {
	AccountID      int    `json:"account_id"`
	CurrentBalance string `json:"current_balance"`
}

type CreateAccountRequest struct {
	AccountID      int    `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}
