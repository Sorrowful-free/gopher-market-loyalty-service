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
		SELECT 
			COALESCE(SUM(CASE WHEN transaction_type = 'ACCRUAL' THEN amount ELSE 0 END), 0) - 
			COALESCE(SUM(CASE WHEN transaction_type = 'WITHDRAWAL' THEN amount ELSE 0 END), 0) as current,
			COALESCE(SUM(CASE WHEN transaction_type = 'WITHDRAWAL' THEN amount ELSE 0 END), 0) as withdrawn
		FROM balance_transactions
		WHERE user_id = $1
	`
	row := r.pgxPool.QueryRow(ctx, query, userID)
	var balanceModel models.BalanceModel
	err := row.Scan(&balanceModel.Current, &balanceModel.Withdrawn)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Если транзакций нет, возвращаем нулевой баланс
			return *models.NewBalanceModel(0, 0), nil
		}
		return models.EmptyBalanceModel, NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to get balance by user id")
	}
	return balanceModel, nil
}
func (r *PGXBalanceRepository) Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error {
	// Используем транзакцию для атомарности операций
	tx, err := r.pgxPool.Begin(ctx)
	if err != nil {
		return NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Вставляем запись о списании в withdrawals
	const insertWithdrawalQuery = `
		INSERT INTO withdrawals (user_id, order_number, sum)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, order_number, sum, processed_at
	`
	row := tx.QueryRow(ctx, insertWithdrawalQuery, userID, orderNumber, sum)
	var withdrawalModel models.WithdrawalModel
	err = row.Scan(&withdrawalModel.ID, &withdrawalModel.UserID, &withdrawalModel.OrderNumber, &withdrawalModel.Sum, &withdrawalModel.ProcessedAt)
	if err != nil {
		return NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to insert withdrawal")
	}

	// Добавляем транзакцию в balance_transactions
	const insertTransactionQuery = `
		INSERT INTO balance_transactions (user_id, order_id, transaction_type, amount)
		VALUES ($1, NULL, 'WITHDRAWAL', $2)
	`
	_, err = tx.Exec(ctx, insertTransactionQuery, userID, sum)
	if err != nil {
		return NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to insert balance transaction")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to commit transaction")
	}

	return nil
}

func (r *PGXBalanceRepository) GetWithdrawals(ctx context.Context, userID int) ([]models.WithdrawalModel, error) {
	const query = `
		SELECT id, user_id, order_number, sum, processed_at
		FROM withdrawals
		WHERE user_id = $1
		ORDER BY processed_at DESC
	`
	rows, err := r.pgxPool.Query(ctx, query, userID)
	if err != nil {
		return []models.WithdrawalModel{}, NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to get withdrawals")
	}
	defer rows.Close()
	var withdrawals = []models.WithdrawalModel{}
	for rows.Next() {
		var withdrawalModel models.WithdrawalModel
		err := rows.Scan(&withdrawalModel.ID, &withdrawalModel.UserID, &withdrawalModel.OrderNumber, &withdrawalModel.Sum, &withdrawalModel.ProcessedAt)
		if err != nil {
			return []models.WithdrawalModel{}, NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to get withdrawals")
		}
		withdrawals = append(withdrawals, withdrawalModel)
	}
	return withdrawals, nil
}

func (r *PGXBalanceRepository) Accrual(ctx context.Context, userID int, orderID int, amount float64) error {
	// Проверяем, не было ли уже начисления за этот заказ
	const checkQuery = `
		SELECT COUNT(*) FROM balance_transactions
		WHERE user_id = $1 AND order_id = $2 AND transaction_type = 'ACCRUAL'
	`
	var count int
	err := r.pgxPool.QueryRow(ctx, checkQuery, userID, orderID).Scan(&count)
	if err != nil {
		return NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to check existing accrual")
	}
	if count > 0 {
		// Начисление уже было, ничего не делаем
		return nil
	}

	const query = `
		INSERT INTO balance_transactions (user_id, order_id, transaction_type, amount)
		VALUES ($1, $2, 'ACCRUAL', $3)
	`
	_, err = r.pgxPool.Exec(ctx, query, userID, orderID, amount)
	if err != nil {
		return NewBalanceRepositoryError(BalanceRepositoryErrorInternalError, "Failed to insert accrual transaction")
	}
	return nil
}
