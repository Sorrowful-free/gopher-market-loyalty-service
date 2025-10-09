package repositories

import (
	"database/sql"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
)

type PGOrderRepository struct {
	db *sql.DB
}

func NewPGOrderRepository(db *sql.DB) OrderRepository {
	return &PGOrderRepository{db: db}
}

func (r *PGOrderRepository) CreateOrder(order models.OrderModel) (models.OrderModel, error) {
	return models.OrderModel{}, nil
}

func (r *PGOrderRepository) GetOrdersList(userID int64) ([]models.OrderModel, error) {
	return []models.OrderModel{}, nil
}

func (r *PGOrderRepository) GetOrder(orderID int64) (models.OrderModel, error) {
	return models.OrderModel{}, nil
}
