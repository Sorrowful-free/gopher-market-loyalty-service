package services

import (
	"context"
	"errors"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/utils"
)

//go:generate mockgen -source=order_service.go -destination=mock_order_service.go -package=services
type OrderService interface {
	CreateOrder(ctx context.Context, userID int, orderNumber string) (models.OrderModel, error)
	GetOrdersList(ctx context.Context, userID int) ([]models.OrderModel, error)
	GetOrder(ctx context.Context, orderNumber string) (models.OrderModel, error)
}

type OrderServiceImpl struct {
	orderRepository repositories.OrderRepository
}

func NewOrderService(orderRepository repositories.OrderRepository) OrderService {
	return &OrderServiceImpl{orderRepository: orderRepository}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, userID int, orderNumber string) (models.OrderModel, error) {
	if ctx.Err() != nil {
		return models.EmptyOrderModel, ctx.Err()
	}

	if !utils.ValidateLuhn(orderNumber) {
		return models.EmptyOrderModel, NewOrderServiceError(OrderServiceErrorOrderIDIsInvalid, "Order id is invalid")
	}

	orderModel, err := s.orderRepository.CreateOrder(ctx, userID, orderNumber)

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
		return models.EmptyOrderModel, err
	}
	return orderModel, nil
}

func (s *OrderServiceImpl) GetOrdersList(ctx context.Context, userID int) ([]models.OrderModel, error) {
	if ctx.Err() != nil {
		return []models.OrderModel{}, ctx.Err()
	}

	orders, err := s.orderRepository.GetOrdersList(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderServiceImpl) GetOrder(ctx context.Context, orderNumber string) (models.OrderModel, error) {

	if ctx.Err() != nil {
		return models.EmptyOrderModel, ctx.Err()
	}

	if !utils.ValidateLuhn(orderNumber) {
		return models.EmptyOrderModel, NewOrderServiceError(OrderServiceErrorOrderIDIsInvalid, "Order id is invalid")
	}

	orderModel, err := s.orderRepository.GetOrder(ctx, orderNumber)
	var orderRepositoryError repositories.OrderRepositoryError
	if errors.As(err, &orderRepositoryError) {
		switch orderRepositoryError.Code {
		case repositories.OrderRepositoryErrorOrderNotFound:
			return models.EmptyOrderModel, NewOrderServiceError(OrderServiceErrorOrderNotFound, "Order not found")
		}
	}
	return orderModel, nil
}
