package repositories

import (
	context "context"

	models "github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
)

//go:generate mockgen -source=balance_repository.go -destination=mock_balance_repository.go -package=repositories
type BalanceRepository interface {
	GetBalance(ctx context.Context, userID string) (models.BalanceModel, error)
	Withdraw(ctx context.Context, userID string, orderID string, sum float64) error
	GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawalModel, error)
}
