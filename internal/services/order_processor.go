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

// OrderProcessor обрабатывает заказы в фоновом режиме
type OrderProcessor struct {
	orderRepository           repositories.OrderRepository
	externalAccrualRepository repositories.ExternalAccrualRepository
	balanceRepository         repositories.BalanceRepository
	logger                    logger.Logger

	// Интервал между проверками заказов
	pollInterval time.Duration
	// Количество заказов для обработки за один цикл
	batchSize int
	// Задержка при получении 429 (rate limit)
	rateLimitDelay time.Duration
}

// NewOrderProcessor создает новый процессор заказов
func NewOrderProcessor(
	orderRepository repositories.OrderRepository,
	externalAccrualRepository repositories.ExternalAccrualRepository,
	balanceRepository repositories.BalanceRepository,
	logger logger.Logger,
) *OrderProcessor {
	return &OrderProcessor{
		orderRepository:           orderRepository,
		externalAccrualRepository:  externalAccrualRepository,
		balanceRepository:         balanceRepository,
		logger:                    logger,
		pollInterval:              5 * time.Second,  // Проверка каждые 5 секунд
		batchSize:                 10,                // Обрабатываем до 10 заказов за раз
		rateLimitDelay:            60 * time.Second,  // Задержка при rate limit
	}
}

// Start запускает фоновую обработку заказов
func (p *OrderProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	// Обрабатываем сразу при старте
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

// processOrders обрабатывает пакет заказов
func (p *OrderProcessor) processOrders(ctx context.Context) {
	// Получаем заказы со статусами NEW и PROCESSING
	orders, err := p.getPendingOrders(ctx)
	if err != nil {
		p.logger.Error("Failed to get pending orders", "error", err)
		return
	}

	if len(orders) == 0 {
		return
	}

	p.logger.Info("Processing orders", "count", len(orders))

	// Обрабатываем заказы пакетами
	for i := 0; i < len(orders) && i < p.batchSize; i++ {
		order := orders[i]
		if err := p.processOrder(ctx, order); err != nil {
			p.logger.Error("Failed to process order",
				"order_number", order.OrderNumber,
				"error", err,
			)
			// Продолжаем обработку других заказов
			continue
		}
	}
}

// getPendingOrders получает заказы со статусами NEW и PROCESSING
func (p *OrderProcessor) getPendingOrders(ctx context.Context) ([]models.OrderModel, error) {
	return p.orderRepository.GetPendingOrders(ctx, p.batchSize)
}

// processOrder обрабатывает один заказ
func (p *OrderProcessor) processOrder(ctx context.Context, order models.OrderModel) error {
	// Получаем информацию о статусе заказа из внешней системы
	scoring, err := p.externalAccrualRepository.GetScoring(ctx, order.OrderNumber)

	var externalAccrualRepositoryError repositories.ExternalAccrualRepositoryError
	if errors.As(err, &externalAccrualRepositoryError) {
		switch externalAccrualRepositoryError.Code {
		case repositories.ExternalAccrualRepositoryErrorOrderTooManyRequests:
			// Rate limit - пропускаем этот заказ, обработаем в следующем цикле
			p.logger.Warn("Rate limit hit for order", "order_number", order.OrderNumber)
			time.Sleep(p.rateLimitDelay)
			return nil

		case repositories.ExternalAccrualRepositoryErrorOrderNotRegistered:
			// Заказ еще не зарегистрирован в системе начислений
			// Это нормально, просто пропускаем
			return nil

		case repositories.ExternalAccrualRepositoryErrorInternalError:
			return fmt.Errorf("external accrual repository error: %w", err)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to get scoring: %w", err)
	}

	// Определяем новый статус и начисление
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
		// Заказ зарегистрирован, но еще не обработан
		newStatus = models.OrderStatusNew
	}

	// Обновляем заказ в БД
	updatedOrder, err := p.orderRepository.UpdateOrder(ctx, order.OrderNumber, newStatus, accrual)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// Если заказ обработан и есть начисление - начисляем баллы на баланс
	if newStatus == models.OrderStatusProcessed && accrual > 0 {
		if err := p.accrualBalance(ctx, updatedOrder); err != nil {
			p.logger.Error("Failed to accrual balance",
				"order_id", updatedOrder.OrderID,
				"user_id", updatedOrder.UserID,
				"accrual", accrual,
				"error", err,
			)
			// Не возвращаем ошибку, т.к. заказ уже обновлен
			// Можно добавить механизм повторных попыток начисления
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

// accrualBalance начисляет баллы на баланс пользователя
func (p *OrderProcessor) accrualBalance(ctx context.Context, order models.OrderModel) error {
	return p.balanceRepository.Accrual(ctx, order.UserID, order.OrderID, order.Accrual)
}

