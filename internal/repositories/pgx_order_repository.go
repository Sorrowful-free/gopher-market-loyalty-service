package repositories

import (
	context "context"
	"errors"
	"time"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXOrderRepository struct {
	pgxPool *pgxpool.Pool
}

func NewPGXOrderRepository(pgxPool *pgxpool.Pool) OrderRepository {
	return &PGXOrderRepository{pgxPool: pgxPool}
}

func (r *PGXOrderRepository) CreateOrder(ctx context.Context, userID int, orderID int) (models.OrderModel, error) {
	const query = `
		INSERT INTO orders (id, user_id, status, accrual, uploaded_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, status, accrual, uploaded_at
	`
	row := r.pgxPool.QueryRow(ctx, query, orderID, userID, models.OrderStatusNew, 0, time.Now().UTC())
	var orderModel models.OrderModel
	err := row.Scan(&orderModel.OrderID, &orderModel.UserID, &orderModel.Status, &orderModel.Accrual, &orderModel.CreatedAt)

	var pgxErr *pgconn.PgError
	if errors.As(err, &pgxErr) {
		if pgxErr.Code == "23505" {
			return models.EMPTY_ORDER_MODEL, NewOrderRepositoryError(OrderRepositoryErrorOrderAlreadyExists, "Order already exists")
		}
	}

	if err != nil {
		return models.EMPTY_ORDER_MODEL, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to create order")
	}
	return orderModel, nil
}

func (r *PGXOrderRepository) GetOrdersList(ctx context.Context, userID int) ([]models.OrderModel, error) {
	const query = `
		SELECT id, user_id, order_id, status, accrual, created_at
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
		SELECT id, user_id, order_id, status, accrual, created_at
		FROM orders
		WHERE order_id = $1
	`
	row := r.pgxPool.QueryRow(ctx, query, orderID)
	var orderModel models.OrderModel
	err := row.Scan(&orderModel.OrderID, &orderModel.UserID, &orderModel.Status, &orderModel.Accrual, &orderModel.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.EMPTY_ORDER_MODEL, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to get order")
		}
		return models.EMPTY_ORDER_MODEL, NewOrderRepositoryError(OrderRepositoryErrorInternalError, "Failed to get order")
	}
	return orderModel, nil
}
