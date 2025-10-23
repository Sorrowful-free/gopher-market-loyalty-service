package repositories

import (
	context "context"
	"database/sql"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
)

type PGUserRepository struct {
	db *sql.DB
}

func NewPGUserRepository(db *sql.DB) UserRepository {
	return &PGUserRepository{db: db}
}

func (r *PGUserRepository) Create(ctx context.Context, login string, password string) (models.UserModel, error) {
	return models.UserModel{}, nil
}

func (r *PGUserRepository) GetByLoginAndPassword(ctx context.Context, login string, password string) (models.UserModel, error) {
	return models.UserModel{}, nil
}

func (r *PGUserRepository) GetBalance(ctx context.Context, userID string) (models.BalanceModel, error) {
	return models.BalanceModel{}, nil
}
