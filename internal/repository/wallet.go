package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/bryanaleron193/wallet-service/internal/model"
	"github.com/bryanaleron193/wallet-service/pkg/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *model.Wallet) error
	GetByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	Withdraw(ctx context.Context, walletID string, amount float64, desc string) (*model.Wallet, string, error)
}

type walletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(ctx context.Context, wallet *model.Wallet) error {
	query := `
		INSERT INTO wallets (user_id, balance)
		VALUES (@userId, @balance)
		RETURNING id, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"userId":  wallet.UserID,
		"balance": wallet.Balance,
	}

	err := r.db.QueryRow(ctx, query, args).Scan(
		&wallet.ID,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("repo.Wallet.Create: %w", err)
	}

	return nil
}

func (r *walletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			balance
		FROM wallets 
		WHERE 
			user_id = $1 
			AND deleted_at IS NULL`

	wallet := &model.Wallet{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.ErrNotFound
		}

		return nil, fmt.Errorf("repo.Wallet.GetByUserID: %w", err)
	}

	return wallet, nil
}

func (r *walletRepository) Withdraw(ctx context.Context, walletID string, amount float64, desc string) (*model.Wallet, string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to start transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	var currentBalance float64

	queryLock := `
		SELECT balance 
		FROM wallets 
		WHERE 
			id = $1 
			AND deleted_at IS NULL 
		FOR UPDATE
	`

	err = tx.QueryRow(ctx, queryLock, walletID).Scan(&currentBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", response.ErrNotFound
		}
		return nil, "", fmt.Errorf("failed to lock wallet row: %w", err)
	}

	if currentBalance < amount {
		return nil, "", response.ErrInsufficient
	}

	newBalance := currentBalance - amount

	queryUpdate := `
		UPDATE wallets 
		SET 
			balance = $1
		WHERE id = $2
	`

	_, err = tx.Exec(ctx, queryUpdate, newBalance, walletID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to update wallet balance: %w", err)
	}

	var transactionID string

	queryInsert := `
		INSERT INTO transactions (wallet_id, amount, transaction_type, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err = tx.QueryRow(ctx, queryInsert, walletID, amount, "withdraw", desc).Scan(&transactionID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to insert transaction record: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &model.Wallet{
		ID:      walletID,
		Balance: newBalance,
	}, transactionID, nil
}
