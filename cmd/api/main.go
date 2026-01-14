package main

import (
	"github.com/bryanaleron193/wallet-service/internal/app"
	"github.com/bryanaleron193/wallet-service/internal/config"
	"github.com/bryanaleron193/wallet-service/pkg/database"
	"github.com/bryanaleron193/wallet-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	_ "github.com/bryanaleron193/wallet-service/docs"
)

//	@title			Wallet API
//	@version		1.0
//	@description	This is a digital wallet service.
//	@host			localhost:8081
//	@BasePath		/

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type 'Bearer TOKEN' to authenticate
func main() {
	cfg := config.Load()
	logger.Setup(cfg.App.Env)

	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer db.Close()
	log.Info().Msg("database connection pool established")

	e := echo.New()

	app.InitMiddleware(e)

	container := app.NewContainer(db, cfg)
	container.RegisterRoutes(e)

	address := ":" + cfg.App.Port
	log.Info().Str("port", cfg.App.Port).Msg("Starting server")

	if err := e.Start(address); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
