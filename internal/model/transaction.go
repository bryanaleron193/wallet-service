package model

import "time"

const (
	TxTypeWithdraw = "withdraw"
	TxTypeDeposit  = "deposit"
)

type Transaction struct {
	ID              string    `json:"id"`
	WalletID        string    `json:"wallet_id"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	Description     *string   `json:"description,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
