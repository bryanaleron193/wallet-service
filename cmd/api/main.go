package main

import (
	"net/http"

	"github.com/bryanaleron193/wallet-service/internal/config"
	"github.com/bryanaleron193/wallet-service/pkg/database"
	"github.com/bryanaleron193/wallet-service/pkg/logger"
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

	address := ":" + cfg.App.Port
	log.Info().Str("port", cfg.App.Port).Msg("Starting server")

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
