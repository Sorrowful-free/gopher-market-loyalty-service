package repositories

import (
	context "context"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
)

//go:generate mockgen -source=order_repository.go -destination=mock_order_repository.go -package=repositories
type OrderRepository interface {
	CreateOrder(ctx context.Context, userID int, order int) (models.OrderModel, error)
	GetOrdersList(ctx context.Context, userID int) ([]models.OrderModel, error)
	GetOrder(ctx context.Context, orderID int) (models.OrderModel, error)
}
