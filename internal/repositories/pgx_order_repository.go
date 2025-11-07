package repositories

import (
	context "context"
	"time"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXOrderRepository struct {
	pgxPool *pgxpool.Pool
}

func NewPGXOrderRepository(pgxPool *pgxpool.Pool) OrderRepository {
	return &PGXOrderRepository{pgxPool: pgxPool}
}

func (r *PGXOrderRepository) CreateOrder(ctx context.Context, userID int, orderID int) (models.OrderModel, error) {
	const selectQuery = `
		SELECT id, user_id
		FROM orders
		WHERE id = $1
	`
	const insertQuery = `
		INSERT INTO orders (id, user_id, status, accrual, uploaded_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, status, accrual, uploaded_at
	`

	tx, err := r.pgxPool.Begin(ctx)
	if err != nil {
		return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to begin transaction")
	}

	row := tx.QueryRow(ctx, selectQuery, orderID)
	var orderIDFromDB int
	var userIDFromDB int
	err = row.Scan(&orderIDFromDB, &userIDFromDB)

	var orderModel models.OrderModel

	if err == pgx.ErrNoRows { // we cannot find order in database
		row = tx.QueryRow(ctx, insertQuery, orderID, userID, models.OrderStatusNew, 0, time.Now().UTC())
		err = row.Scan(&orderModel.OrderID, &orderModel.UserID, &orderModel.Status, &orderModel.Accrual, &orderModel.CreatedAt)
		if err != nil {
			return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to insert order")
		}

		err = tx.Commit(ctx)
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to rollback transaction")
			}
			return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to commit transaction")
		}
		return orderModel, nil
	}

	if userIDFromDB != userID {
		return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorOrderCreatedOtherUser, "Order created by other user")
	} else {
		return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorOrderAlreadyExists, "Order already exists")
	}
}

func (r *PGXOrderRepository) GetOrdersList(ctx context.Context, userID int) ([]models.OrderModel, error) {
	const query = `
		SELECT id, user_id, status, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1
	`
	rows, err := r.pgxPool.Query(ctx, query, userID)
	if err != nil {
		return []models.OrderModel{}, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to get orders list")
	}
	defer rows.Close()
	var orders = []models.OrderModel{}
	for rows.Next() {
		var orderModel models.OrderModel
		err := rows.Scan(&orderModel.OrderID, &orderModel.UserID, &orderModel.Status, &orderModel.Accrual, &orderModel.CreatedAt)
		if err != nil {
			return []models.OrderModel{}, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to get orders list")
		}
		orders = append(orders, orderModel)
	}
	return orders, nil
}

func (r *PGXOrderRepository) GetOrder(ctx context.Context, orderID int) (models.OrderModel, error) {
	const query = `
		SELECT id, user_id, status, accrual, uploaded_at
		FROM orders
		WHERE id = $1
	`
	row := r.pgxPool.QueryRow(ctx, query, orderID)
	var orderModel models.OrderModel
	err := row.Scan(&orderModel.OrderID, &orderModel.UserID, &orderModel.Status, &orderModel.Accrual, &orderModel.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to get order")
		}
		return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to get order")
	}
	return orderModel, nil
}
