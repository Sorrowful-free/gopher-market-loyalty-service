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

func (r *PGXOrderRepository) CreateOrder(ctx context.Context, userID int, orderNumber string, orderStatus models.OrderStatus, accrual float64) (models.OrderModel, error) {
	const selectQuery = `
		SELECT order_number, user_id
		FROM orders
		WHERE order_number = $1
	`
	const insertQuery = `
		INSERT INTO orders (order_number, user_id, status, accrual, uploaded_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, order_number, user_id, status, accrual, uploaded_at
	`

	tx, err := r.pgxPool.Begin(ctx)
	if err != nil {
		return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to begin transaction")
	}

	row := tx.QueryRow(ctx, selectQuery, orderNumber)
	var orderNumberFromDB string
	var userIDFromDB int
	err = row.Scan(&orderNumberFromDB, &userIDFromDB)

	var orderModel models.OrderModel

	if err == pgx.ErrNoRows { // we cannot find order in database
		row = tx.QueryRow(ctx, insertQuery, orderNumber, userID, orderStatus, accrual, time.Now().UTC())
		err = row.Scan(&orderModel.OrderID, &orderModel.OrderNumber, &orderModel.UserID, &orderModel.Status, &orderModel.Accrual, &orderModel.UploadedAt)
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
		SELECT id, order_number, user_id, status, accrual, uploaded_at
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
		err := rows.Scan(&orderModel.OrderID, &orderModel.OrderNumber, &orderModel.UserID, &orderModel.Status, &orderModel.Accrual, &orderModel.UploadedAt)
		if err != nil {
			return []models.OrderModel{}, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to get orders list")
		}
		orders = append(orders, orderModel)
	}
	return orders, nil
}

func (r *PGXOrderRepository) GetOrder(ctx context.Context, orderNumber string) (models.OrderModel, error) {
	const query = `
		SELECT id, user_id, order_number, status, accrual, uploaded_at
		FROM orders
		WHERE order_number = $1
	`
	row := r.pgxPool.QueryRow(ctx, query, orderNumber)
	var orderModel models.OrderModel
	err := row.Scan(&orderModel.OrderID, &orderModel.UserID, &orderModel.OrderNumber, &orderModel.Status, &orderModel.Accrual, &orderModel.UploadedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to get order")
		}
		return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to get order")
	}
	return orderModel, nil
}

func (r *PGXOrderRepository) UpdateOrder(ctx context.Context, orderNumber string, orderStatus models.OrderStatus, accrual float64) (models.OrderModel, error) {
	const query = `
		UPDATE orders
		SET status = $1, accrual = $2, uploaded_at = $3
		WHERE order_number = $4
		RETURNING id, order_number, user_id, status, accrual, uploaded_at
	`
	row := r.pgxPool.QueryRow(ctx, query, orderStatus, accrual, time.Now().UTC(), orderNumber)
	var orderModel models.OrderModel
	err := row.Scan(&orderModel.OrderID, &orderModel.OrderNumber, &orderModel.UserID, &orderModel.Status, &orderModel.Accrual, time.Now().UTC())
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorOrderNotFound, "Order not found")
		}
		return models.EmptyOrderModel, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to update order")
	}
	return orderModel, nil
}
