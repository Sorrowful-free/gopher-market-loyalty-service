package services

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
)

type OrderService interface {
	CreateOrder(userID string, order string) (models.OrderModel, error)
	GetOrdersList(userID string) ([]models.OrderModel, error)
	GetOrder(orderID string) (models.OrderModel, error)
}

type OrderServiceImpl struct {
	orderRepository repositories.OrderRepository
}

func NewOrderService(orderRepository repositories.OrderRepository) OrderService {
	return &OrderServiceImpl{orderRepository: orderRepository}
}

func (s *OrderServiceImpl) CreateOrder(userID string, order string) (models.OrderModel, error) {
	return s.orderRepository.CreateOrder(userID, order)
}

func (s *OrderServiceImpl) GetOrdersList(userID string) ([]models.OrderModel, error) {
	return s.orderRepository.GetOrdersList(userID)
}

func (s *OrderServiceImpl) GetOrder(orderID string) (models.OrderModel, error) {
	return s.orderRepository.GetOrder(orderID)
}
