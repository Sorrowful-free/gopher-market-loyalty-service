package repositories

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXOrderRepository struct {
	pgxPool *pgxpool.Pool
}

func NewPGXOrderRepository(pgxPool *pgxpool.Pool) OrderRepository {
	return &PGXOrderRepository{pgxPool: pgxPool}
}

func (r *PGXOrderRepository) CreateOrder(userID string, order string) (models.OrderModel, error) {
	return models.OrderModel{}, nil
}

func (r *PGXOrderRepository) GetOrdersList(userID string) ([]models.OrderModel, error) {
	return []models.OrderModel{}, nil
}

func (r *PGXOrderRepository) GetOrder(orderID string) (models.OrderModel, error) {
	return models.OrderModel{}, nil
}
