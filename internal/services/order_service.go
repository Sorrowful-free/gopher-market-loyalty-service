package services

import (
	"context"
	"errors"
	"fmt"

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
	orderRepository           repositories.OrderRepository
	externalAccrualRepository repositories.ExternalAccrualRepository
}

func NewOrderService(orderRepository repositories.OrderRepository, externalAccrualRepository repositories.ExternalAccrualRepository) OrderService {
	return &OrderServiceImpl{orderRepository: orderRepository, externalAccrualRepository: externalAccrualRepository}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, userID int, orderNumber string) (models.OrderModel, error) {
	if ctx.Err() != nil {
		return models.EmptyOrderModel, ctx.Err()
	}

	if !utils.ValidateLuhn(orderNumber) {
		return models.EmptyOrderModel, NewOrderServiceError(OrderServiceErrorOrderIDIsInvalid, "Order id is invalid")
	}

	orderModel, err := s.orderRepository.CreateOrder(ctx, userID, orderNumber, models.OrderStatusNew, 0)

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
		return models.EmptyOrderModel, fmt.Errorf("failed to create order: %w", err)
	}

	scoring, err := s.externalAccrualRepository.GetScoring(ctx, orderNumber)

	var externalAccrualRepositoryError repositories.ExternalAccrualRepositoryError
	if errors.As(err, &externalAccrualRepositoryError) {
		switch externalAccrualRepositoryError.Code {
		case repositories.ExternalAccrualRepositoryErrorOrderTooManyRequests, repositories.ExternalAccrualRepositoryErrorOrderNotRegistered: // that means we need to wait for the order to be processed
			return orderModel, nil // that expected behavior
		case repositories.ExternalAccrualRepositoryErrorInternalError:
			return models.EmptyOrderModel, NewOrderServiceError(OrderServiceErrorInternalError, fmt.Sprintf("External accrual repository error: %s", externalAccrualRepositoryError.Message))
		}
	}

	if err != nil {
		return models.EmptyOrderModel, fmt.Errorf("failed to get scoring: %w", err)
	}

	var orderStatus models.OrderStatus = models.OrderStatusNew
	var accrual float64 = 0
	switch scoring.Status {
	case models.ScoringStatusProcessed:
		orderStatus = models.OrderStatusProcessed
		accrual = scoring.Accrual
	case models.ScoringStatusInvalid:
		orderStatus = models.OrderStatusInvalid
	case models.ScoringStatusProcessing:
		orderStatus = models.OrderStatusProcessing
	case models.ScoringStatusRegistered:
		orderStatus = models.OrderStatusNew
	}

	orderModel, err = s.orderRepository.UpdateOrder(ctx, orderNumber, orderStatus, accrual)
	if err != nil {
		return models.EmptyOrderModel, fmt.Errorf("failed to update order: %w", err)
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

	for idx, order := range orders {

		if order.Status == models.OrderStatusProcessed ||
			order.Status == models.OrderStatusInvalid { //||
			// order.UploadedAt.Before(time.Now().UTC().Add(-1*time.Second*15)) { // 15 seconds is the time to wait for the order to be processed
			continue
		}

		scoring, err := s.externalAccrualRepository.GetScoring(ctx, order.OrderNumber)

		var externalAccrualRepositoryError repositories.ExternalAccrualRepositoryError
		if errors.As(err, &externalAccrualRepositoryError) {
			switch externalAccrualRepositoryError.Code {
			case repositories.ExternalAccrualRepositoryErrorOrderTooManyRequests, repositories.ExternalAccrualRepositoryErrorOrderNotRegistered:
				_, err = s.orderRepository.UpdateOrder(ctx, order.OrderNumber, order.Status, order.Accrual) // refresh order in case of order not registered
				if err != nil {
					return nil, fmt.Errorf("failed to update order: %w", err)
				}
				continue
			}
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get scoring: %w", err)
		}
		if scoring.Status == models.ScoringStatusProcessed {
			order.Status = models.OrderStatusProcessed
			order.Accrual = scoring.Accrual
		}
		if scoring.Status == models.ScoringStatusInvalid {
			order.Status = models.OrderStatusInvalid
		}
		if scoring.Status == models.ScoringStatusProcessing {
			order.Status = models.OrderStatusProcessing
		}
		order, err = s.orderRepository.UpdateOrder(ctx, order.OrderNumber, order.Status, order.Accrual)

		var orderRepositoryError repositories.OrderRepositoryError
		if errors.As(err, &orderRepositoryError) {
			switch orderRepositoryError.Code {
			case repositories.OrderRepositoryErrorOrderNotFound:
				return nil, NewOrderServiceError(OrderServiceErrorOrderNotFound, "Order not found")
			}
		}

		if err != nil {
			return nil, err
		}
		orders[idx] = order
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
