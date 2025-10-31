package repositories

import (
	context "context"

	models "github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXBalanceRepository struct {
	pgxPool *pgxpool.Pool
}

func NewPGXBalanceRepository(pgxPool *pgxpool.Pool) *PGXBalanceRepository {
	return &PGXBalanceRepository{
		pgxPool: pgxPool,
	}
}

func (r *PGXBalanceRepository) GetBalance(ctx context.Context, userID int) (models.BalanceModel, error) {
	return models.EMPTY_BALANCE_MODEL, nil
}
func (r *PGXBalanceRepository) Withdraw(ctx context.Context, userID int, orderID int, sum float64) error {
	return nil
}
func (r *PGXBalanceRepository) GetWithdrawals(ctx context.Context, userID int) ([]models.WithdrawalModel, error) {
	return []models.WithdrawalModel{}, nil
}
