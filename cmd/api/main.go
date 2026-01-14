package main

import (
	"github.com/bryanaleron193/wallet-service/internal/app"
	"github.com/bryanaleron193/wallet-service/internal/config"
	"github.com/bryanaleron193/wallet-service/pkg/database"
	"github.com/bryanaleron193/wallet-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

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

	container := app.NewContainer(db)
	container.RegisterRoutes(e)

	address := ":" + cfg.App.Port
	log.Info().Str("port", cfg.App.Port).Msg("Starting server")

	if err := e.Start(address); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
