package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
)

type OrderProcessor struct {
	orderRepository           repositories.OrderRepository
	externalAccrualRepository repositories.ExternalAccrualRepository
	balanceRepository         repositories.BalanceRepository
	logger                    logger.Logger

	pollInterval   time.Duration
	batchSize      int
	rateLimitDelay time.Duration
}

func NewOrderProcessor(
	orderRepository repositories.OrderRepository,
	externalAccrualRepository repositories.ExternalAccrualRepository,
	balanceRepository repositories.BalanceRepository,
	logger logger.Logger,
) *OrderProcessor {
	return &OrderProcessor{
		orderRepository:           orderRepository,
		externalAccrualRepository: externalAccrualRepository,
		balanceRepository:         balanceRepository,
		logger:                    logger,
		pollInterval:              1 * time.Second,
		batchSize:                 10,
		rateLimitDelay:            1 * time.Second,
	}
}

func (p *OrderProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	p.processOrders(ctx)

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Order processor stopped", "reason", ctx.Err())
			return
		case <-ticker.C:
			p.processOrders(ctx)
		}
	}
}

func (p *OrderProcessor) processOrders(ctx context.Context) {
	orders, err := p.getPendingOrders(ctx)
	if err != nil {
		p.logger.Error("Failed to get pending orders", "error", err)
		return
	}

	if len(orders) == 0 {
		return
	}

	p.logger.Info("Processing orders", "count", len(orders))

	for i := 0; i < len(orders) && i < p.batchSize; i++ {
		order := orders[i]
		if err := p.processOrder(ctx, order); err != nil {
			p.logger.Error("Failed to process order",
				"order_number", order.OrderNumber,
				"error", err,
			)
			continue
		}
	}
}

func (p *OrderProcessor) getPendingOrders(ctx context.Context) ([]models.OrderModel, error) {
	return p.orderRepository.GetPendingOrders(ctx, p.batchSize)
}

func (p *OrderProcessor) processOrder(ctx context.Context, order models.OrderModel) error {
	scoring, err := p.externalAccrualRepository.GetScoring(ctx, order.OrderNumber)

	var externalAccrualRepositoryError repositories.ExternalAccrualRepositoryError
	if errors.As(err, &externalAccrualRepositoryError) {
		switch externalAccrualRepositoryError.Code {
		case repositories.ExternalAccrualRepositoryErrorOrderTooManyRequests:
			p.logger.Warn("Rate limit hit for order", "order_number", order.OrderNumber)
			time.Sleep(p.rateLimitDelay)
			return nil

		case repositories.ExternalAccrualRepositoryErrorOrderNotRegistered:
			return nil

		case repositories.ExternalAccrualRepositoryErrorInternalError:
			return fmt.Errorf("external accrual repository error: %w", err)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to get scoring: %w", err)
	}

	var newStatus models.OrderStatus = order.Status
	var accrual float64 = order.Accrual

	switch scoring.Status {
	case models.ScoringStatusProcessed:
		newStatus = models.OrderStatusProcessed
		accrual = scoring.Accrual

	case models.ScoringStatusInvalid:
		newStatus = models.OrderStatusInvalid

	case models.ScoringStatusProcessing:
		newStatus = models.OrderStatusProcessing

	case models.ScoringStatusRegistered:
		newStatus = models.OrderStatusNew
	}

	updatedOrder, err := p.orderRepository.UpdateOrder(ctx, order.OrderNumber, newStatus, accrual)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if newStatus == models.OrderStatusProcessed && accrual > 0 {
		if err := p.accrualBalance(ctx, updatedOrder); err != nil {
			p.logger.Error("Failed to accrual balance",
				"order_id", updatedOrder.OrderID,
				"user_id", updatedOrder.UserID,
				"accrual", accrual,
				"error", err,
			)
		} else {
			p.logger.Info("Balance accrued",
				"order_id", updatedOrder.OrderID,
				"user_id", updatedOrder.UserID,
				"accrual", accrual,
			)
		}
	}

	return nil
}

func (p *OrderProcessor) accrualBalance(ctx context.Context, order models.OrderModel) error {
	return p.balanceRepository.Accrual(ctx, order.UserID, order.OrderID, order.Accrual)
}
