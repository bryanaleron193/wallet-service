package repository

import (
	"context"
	"fmt"

	"github.com/bryanaleron193/wallet-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *model.Transaction) error
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *model.Transaction) error {
	query := `
		INSERT INTO transactions (wallet_id, amount, transaction_type, description)
		VALUES (@walletId, @amount, @txType, @description)
		RETURNING id, created_at
	`

	args := pgx.NamedArgs{
		"walletId":    tx.WalletID,
		"amount":      tx.Amount,
		"txType":      tx.TransactionType,
		"description": tx.Description,
	}

	err := r.db.QueryRow(ctx, query, args).Scan(&tx.ID, &tx.CreatedAt)
	if err != nil {
		return fmt.Errorf("repo.Transaction.Create: %w", err)
	}

	return nil
}
