package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/bryanaleron193/wallet-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, full_name, email)
		VALUES(@username, @fullName, @email)
		RETURNING id, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"username": user.Username,
		"fullName": user.FullName,
		"email":    user.Email,
	}

	err := r.db.QueryRow(
		ctx,
		query,
		args,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("repo.User.Create: %w", err)
	}

	return nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT 
			id, 
			username, 
			email, 
			created_at 
		FROM users 
		WHERE 
			username = $1 
		LIMIT 1`

	user := &model.User{}

	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
