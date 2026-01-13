package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bryanaleron193/wallet-service/internal/config"
	"github.com/bryanaleron193/wallet-service/internal/model"
	"github.com/bryanaleron193/wallet-service/internal/repository"
	"github.com/bryanaleron193/wallet-service/internal/service"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var baseNames = []string{
	"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi",
	"ivan", "judy", "karl", "laura", "mallory", "nina", "oscar", "peggy",
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	logger := log.With().Str("module", "seeder").Logger()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg := config.Load()

	db, err := pgxpool.New(ctx, cfg.DB.URL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to connect to database")
	}
	defer db.Close()

	userService := service.NewUserService(repository.NewUserRepository(db))
	walletService := service.NewWalletService(repository.NewWalletRepository(db))
	transactionService := service.NewTransactionService(repository.NewTransactionRepository(db))

	num := 20
	logger.Info().Int("count", num).Msg("Starting user seeding...")

	users := seedUsers(ctx, userService, num, logger)
	wallets := seedWallets(ctx, walletService, users, logger)
	seedTransactions(ctx, transactionService, wallets, logger)

	logger.Info().Msg("Seeding process finished!")
}

func seedUsers(ctx context.Context, svc service.UserService, count int, logger zerolog.Logger) []*model.User {
	var createdUsers []*model.User

	for i := 0; i < count; i++ {
		base := baseNames[i%len(baseNames)]
		username := fmt.Sprintf("%s_%d", base, i)
		fullName := fmt.Sprintf("%s %d", capitalize(base), i)
		email := fmt.Sprintf("%s@example.com", username)

		user, err := svc.CreateUser(ctx, username, fullName, email)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				logger.Warn().Msg("User already exists, skipping...")
			} else {
				logger.Error().Err(err).Msg("Actual database failure!")
			}
			continue
		}

		logger.Info().
			Str("id", user.ID).
			Str("username", user.Username).
			Msg("Seeded user success")

		createdUsers = append(createdUsers, user)
	}

	return createdUsers
}

func seedWallets(ctx context.Context, svc service.WalletService, users []*model.User, logger zerolog.Logger) []*model.Wallet {
	var createdWallets []*model.Wallet
	initialBalance := 1000000.0

	for _, user := range users {
		wallet, err := svc.CreateWallet(ctx, user.ID, initialBalance)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				logger.Warn().Str("user_id", user.ID).Msg("Wallet already exists, skipping...")
				continue
			}
			logger.Error().Err(err).Str("user_id", user.ID).Msg("Failed to create wallet")
			continue
		}

		logger.Info().
			Str("user_id", user.ID).
			Str("balance", strconv.FormatFloat(initialBalance, 'f', 2, 64)).
			Msg("Seeded Wallet created")

		createdWallets = append(createdWallets, wallet)
	}

	return createdWallets
}

func seedTransactions(ctx context.Context, svc service.TransactionService, wallets []*model.Wallet, logger zerolog.Logger) {
	initialDescription := "Initial Deposit"

	for _, wallet := range wallets {
		_, err := svc.CreateTransaction(
			ctx,
			wallet.ID,
			wallet.Balance,
			model.TxTypeDeposit,
			&initialDescription,
		)

		if err != nil {
			logger.Error().
				Err(err).
				Str("wallet_id", wallet.ID).
				Msg("Failed to seed initial transaction log")
			continue
		}

		logger.Info().Str("wallet_id", wallet.ID).Msg("Initial transaction recorded")
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	return fmt.Sprintf("%c%s", s[0]-32, s[1:])
}
