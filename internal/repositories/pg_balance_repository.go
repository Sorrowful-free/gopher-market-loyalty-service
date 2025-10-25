package repositories

import (
	context "context"
	"database/sql"

	models "github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
)

type PGBalanceRepository struct {
	db *sql.DB
}

func NewPGBalanceRepository(db *sql.DB) *PGBalanceRepository {
	return &PGBalanceRepository{
		db: db,
	}
}

func (r *PGBalanceRepository) GetBalance(ctx context.Context, userID string) (models.BalanceModel, error) {
	return models.EMPTY_BALANCE_MODEL, nil
}
func (r *PGBalanceRepository) Withdraw(ctx context.Context, userID string, orderID string, sum float64) error {
	return nil
}
func (r *PGBalanceRepository) GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawalModel, error) {
	return []models.WithdrawalModel{}, nil
}
