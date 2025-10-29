package repositories

import (
	context "context"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXUserRepository struct {
	pgxPool *pgxpool.Pool
}

func NewPGXUserRepository(pgxPool *pgxpool.Pool) UserRepository {
	return &PGXUserRepository{pgxPool: pgxPool}
}

func (r *PGXUserRepository) Create(ctx context.Context, login string, password string) (models.UserModel, error) {
	return models.UserModel{}, nil
}

func (r *PGXUserRepository) GetByLoginAndPassword(ctx context.Context, login string, password string) (models.UserModel, error) {
	return models.UserModel{}, nil
}

func (r *PGXUserRepository) GetBalance(ctx context.Context, userID string) (models.BalanceModel, error) {
	return models.BalanceModel{}, nil
}
