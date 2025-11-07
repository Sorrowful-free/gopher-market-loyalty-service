package repositories

import (
	context "context"

	models "github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/jackc/pgx/v5"
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
	const query = `
		SELECT current, withdrawn
		FROM balance
		WHERE user_id = $1
	`
	row := r.pgxPool.QueryRow(ctx, query, userID)
	var balanceModel models.BalanceModel
	err := row.Scan(&balanceModel.Current, &balanceModel.Withdrawn)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.EmptyBalanceModel, NewBalanceRepositoryError(BalanceRepositoryErrorUserNotFound, "User not found")
		}
		return models.EmptyBalanceModel, NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to get balance by user id")
	}
	return balanceModel, nil
}
func (r *PGXBalanceRepository) Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error {
	const query = `
		INSERT INTO withdrawals (user_id, order_number, sum)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, order_number, sum, processed_at
	`
	row := r.pgxPool.QueryRow(ctx, query, userID, orderNumber, sum)
	var withdrawalModel models.WithdrawalModel
	err := row.Scan(&withdrawalModel.ID, &withdrawalModel.UserID, &withdrawalModel.OrderID, &withdrawalModel.Sum, &withdrawalModel.ProcessedAt)
	if err != nil {
		return NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to withdraw")
	}
	return nil
}

func (r *PGXBalanceRepository) GetWithdrawals(ctx context.Context, userID int) ([]models.WithdrawalModel, error) {
	const query = `
		SELECT id, user_id, order_id, sum, processed_at
		FROM withdrawals
		WHERE user_id = $1
	`
	rows, err := r.pgxPool.Query(ctx, query, userID)
	if err != nil {
		return []models.WithdrawalModel{}, NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to get withdrawals")
	}
	defer rows.Close()
	var withdrawals = []models.WithdrawalModel{}
	for rows.Next() {
		var withdrawalModel models.WithdrawalModel
		err := rows.Scan(&withdrawalModel.ID, &withdrawalModel.UserID, &withdrawalModel.OrderID, &withdrawalModel.Sum, &withdrawalModel.ProcessedAt)
		if err != nil {
			return []models.WithdrawalModel{}, NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to get withdrawals")
		}
		withdrawals = append(withdrawals, withdrawalModel)
	}
	return withdrawals, nil
}
