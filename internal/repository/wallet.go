package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/bryanaleron193/wallet-service/internal/model"
	"github.com/bryanaleron193/wallet-service/pkg/apperror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *model.Wallet) error
	GetByUserID(ctx context.Context, userID string) (*model.Wallet, error)
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
		RETURNING id, version, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"userId":  wallet.UserID,
		"balance": wallet.Balance,
	}

	err := r.db.QueryRow(ctx, query, args).Scan(
		&wallet.ID,
		&wallet.Version,
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
			return nil, apperror.ErrNotFound
		}

		return nil, fmt.Errorf("repo.Wallet.GetByUserID: %w", err)
	}

	return wallet, nil
}
