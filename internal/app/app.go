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
	HealthHandler *handler.HealthHandler
	UserHandler   *handler.UserHandler
	WalletHandler *handler.WalletHandler
}

func NewContainer(db *pgxpool.Pool, cfg *config.Config) *Container {
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	userService := service.NewUserService(userRepo, cfg)
	walletService := service.NewWalletService(walletRepo)

	return &Container{
		Config:        cfg,
		HealthHandler: handler.NewHealthHandler(),
		UserHandler:   handler.NewUserHandler(userService),
		WalletHandler: handler.NewWalletHandler(walletService),
	}
}
