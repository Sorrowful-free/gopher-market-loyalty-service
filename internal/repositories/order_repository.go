package repositories

import "github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"

type OrderRepository interface {
	CreateOrder(order models.OrderModel) (models.OrderModel, error)
	GetOrdersList(userID int64) ([]models.OrderModel, error)
	GetOrder(orderID int64) (models.OrderModel, error)
}
