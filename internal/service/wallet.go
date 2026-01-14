package service

import (
	"context"
	"fmt"

	"github.com/bryanaleron193/wallet-service/internal/model"
	"github.com/bryanaleron193/wallet-service/internal/repository"
	"github.com/bryanaleron193/wallet-service/pkg/apperror"
	"github.com/bryanaleron193/wallet-service/pkg/utils"
)

type WalletService interface {
	CreateWallet(ctx context.Context, userID string, amount float64) (*model.Wallet, error)
	GetByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	Withdraw(ctx context.Context, userID string, amount float64, desc string) (*model.Wallet, string, error)
}

type walletService struct {
	walletRepo repository.WalletRepository
}

func NewWalletService(wr repository.WalletRepository) WalletService {
	return &walletService{
		walletRepo: wr,
	}
}

func (s *walletService) CreateWallet(ctx context.Context, userID string, amount float64) (*model.Wallet, error) {
	if amount < 0 {
		return nil, fmt.Errorf("initial balance cannot be negative: %w", apperror.ErrInvalidInput)
	}

	wallet := &model.Wallet{
		UserID:  userID,
		Balance: amount,
	}

	err := s.walletRepo.Create(ctx, wallet)
	if err != nil {
		return nil, fmt.Errorf("service.CreateInitialWallet: %w", err)
	}

	return wallet, nil
}

func (s *walletService) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required: %w", apperror.ErrInvalidInput)
	}

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service.GetByUserID: %w", err)
	}

	return wallet, nil
}

func (s *walletService) Withdraw(ctx context.Context, userID string, amount float64, desc string) (*model.Wallet, string, error) {
	if amount <= 0 {
		return nil, "", fmt.Errorf("amount must be greater than zero: %w", apperror.ErrInvalidInput)
	}

	if desc == "" {
		desc = fmt.Sprintf("Withdrawal of %s", utils.FormatRupiah(amount))
	}

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, "", fmt.Errorf("service.Withdraw: %w", err)
	}

	updatedWallet, transactionID, err := s.walletRepo.Withdraw(ctx, wallet.ID, amount, desc)
	if err != nil {
		return nil, "", fmt.Errorf("service.Withdraw execution: %w", err)
	}

	return updatedWallet, transactionID, nil
}
