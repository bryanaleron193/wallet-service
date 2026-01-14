package app

import (
	"github.com/bryanaleron193/wallet-service/internal/handler"
	"github.com/bryanaleron193/wallet-service/internal/repository"
	"github.com/bryanaleron193/wallet-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	WalletHandler *handler.WalletHandler
}

func NewContainer(db *pgxpool.Pool) *Container {
	walletRepo := repository.NewWalletRepository(db)

	walletService := service.NewWalletService(walletRepo)

	return &Container{
		WalletHandler: handler.NewWalletHandler(walletService),
	}
}
