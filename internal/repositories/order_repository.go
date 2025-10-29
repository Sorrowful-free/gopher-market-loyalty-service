package repositories

import "github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"

//go:generate mockgen -source=order_repository.go -destination=mock_order_repository.go -package=repositories
type OrderRepository interface {
	CreateOrder(userID int, order int) (models.OrderModel, error)
	GetOrdersList(userID int) ([]models.OrderModel, error)
	GetOrder(orderID int) (models.OrderModel, error)
}
