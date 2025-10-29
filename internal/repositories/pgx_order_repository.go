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

func (r *PGXOrderRepository) CreateOrder(userID int, order int) (models.OrderModel, error) {
	return models.OrderModel{}, nil
}

func (r *PGXOrderRepository) GetOrdersList(userID int) ([]models.OrderModel, error) {
	return []models.OrderModel{}, nil
}

func (r *PGXOrderRepository) GetOrder(orderID int) (models.OrderModel, error) {
	return models.OrderModel{}, nil
}
