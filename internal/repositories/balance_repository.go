package repositories

import (
	context "context"

	models "github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
)

//go:generate mockgen -source=balance_repository.go -destination=mock_balance_repository.go -package=repositories
type BalanceRepository interface {
	GetBalance(ctx context.Context, userID int) (models.BalanceModel, error)
	Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error
	GetWithdrawals(ctx context.Context, userID int) ([]models.WithdrawalModel, error)
	Accrual(ctx context.Context, userID int, orderID int, amount float64) error
}
