package app

import (
	"github.com/bryanaleron193/wallet-service/internal/config"
	"github.com/bryanaleron193/wallet-service/internal/handler"
	"github.com/bryanaleron193/wallet-service/internal/repository"
	"github.com/bryanaleron193/wallet-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Config        *config.Config
	WalletHandler *handler.WalletHandler
	UserHandler   *handler.UserHandler
}

func NewContainer(db *pgxpool.Pool, cfg *config.Config) *Container {
	walletRepo := repository.NewWalletRepository(db)
	userRepo := repository.NewUserRepository(db)

	walletService := service.NewWalletService(walletRepo)
	userService := service.NewUserService(userRepo, cfg)

	return &Container{
		Config:        cfg,
		WalletHandler: handler.NewWalletHandler(walletService),
		UserHandler:   handler.NewUserHandler(userService),
	}
}
