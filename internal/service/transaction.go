package service

import (
	"context"
	"fmt"

	"github.com/bryanaleron193/wallet-service/internal/model"
	"github.com/bryanaleron193/wallet-service/internal/repository"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, walletID string, amount float64, txType string, description *string) (*model.Transaction, error)
}

type transactionService struct {
	txRepo repository.TransactionRepository
}

func NewTransactionService(tr repository.TransactionRepository) TransactionService {
	return &transactionService{
		txRepo: tr,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, walletID string, amount float64, txType string, description *string) (*model.Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("transaction amount must be positive")
	}

	tx := &model.Transaction{
		WalletID:        walletID,
		Amount:          amount,
		TransactionType: txType,
		Description:     description,
	}

	err := s.txRepo.Create(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("service.CreateTransaction: %w", err)
	}

	return tx, nil
}
