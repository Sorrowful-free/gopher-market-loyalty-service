package services

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
)

type OrderService interface {
	CreateOrder(order models.OrderModel) (models.OrderModel, error)
	GetOrdersList(userID int64) ([]models.OrderModel, error)
	GetOrder(orderID int64) (models.OrderModel, error)
}

type OrderServiceImpl struct {
	orderRepository repositories.OrderRepository
}

func NewOrderService(orderRepository repositories.OrderRepository) OrderService {
	return &OrderServiceImpl{orderRepository: orderRepository}
}

func (s *OrderServiceImpl) CreateOrder(order models.OrderModel) (models.OrderModel, error) {
	return s.orderRepository.CreateOrder(order)
}

func (s *OrderServiceImpl) GetOrdersList(userID int64) ([]models.OrderModel, error) {
	return s.orderRepository.GetOrdersList(userID)
}

func (s *OrderServiceImpl) GetOrder(orderID int64) (models.OrderModel, error) {
	return s.orderRepository.GetOrder(orderID)
}
