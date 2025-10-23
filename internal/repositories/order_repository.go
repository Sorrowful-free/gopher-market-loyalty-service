package repositories

import "github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"

//go:generate mockgen -source=order_repository.go -destination=mock_order_repository.go -package=repositories
type OrderRepository interface {
	CreateOrder(userID string, order string) (models.OrderModel, error)
	GetOrdersList(userID string) ([]models.OrderModel, error)
	GetOrder(orderID string) (models.OrderModel, error)
}
