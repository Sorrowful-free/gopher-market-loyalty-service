package services

import (
	"errors"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/utils"
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
	if !utils.ValidateLuhn(order) {
		return models.EMPTY_ORDER_MODEL, NewOrderServiceError(OrderServiceErrorOrderIdIsInvalid, "Order id is invalid")
	}

	orderModel, err := s.orderRepository.CreateOrder(userID, order)

	var orderRepositoryError repositories.OrderRepositoryError
	if errors.As(err, &orderRepositoryError) {
		switch orderRepositoryError.Code {

		case repositories.OrderRepositoryErrorOrderAlreadyExists:
			return models.OrderModel{}, NewOrderServiceError(OrderServiceErrorOrderAlreadyExists, "Order already exists")
		case repositories.OrderRepositoryErrorOrderCreatedOtherUser:
			return models.OrderModel{}, NewOrderServiceError(OrderServiceErrorOrderCreatedOtherUser, "Order created by other user")
		}
	}
	if err != nil {
		return models.EMPTY_ORDER_MODEL, err
	}
	return orderModel, nil
}

func (s *OrderServiceImpl) GetOrdersList(userID string) ([]models.OrderModel, error) {
	orders, err := s.orderRepository.GetOrdersList(userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderServiceImpl) GetOrder(orderID string) (models.OrderModel, error) {

	if !utils.ValidateLuhn(orderID) {
		return models.EMPTY_ORDER_MODEL, NewOrderServiceError(OrderServiceErrorOrderIdIsInvalid, "Order id is invalid")
	}

	orderModel, err := s.orderRepository.GetOrder(orderID)
	var orderRepositoryError repositories.OrderRepositoryError
	if errors.As(err, &orderRepositoryError) {
		switch orderRepositoryError.Code {
		case repositories.OrderRepositoryErrorOrderNotFound:
			return models.EMPTY_ORDER_MODEL, NewOrderServiceError(OrderServiceErrorOrderNotFound, "Order not found")
		}
	}
	return orderModel, nil
}
