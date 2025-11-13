package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
)

type OrderProcessorPool struct {
	orderRepository           repositories.OrderRepository
	externalAccrualRepository repositories.ExternalAccrualRepository
	balanceRepository         repositories.BalanceRepository
	logger                    logger.Logger

	orderQueue chan models.OrderModel

	workerCount int

	pollInterval time.Duration

	batchSize int

	rateLimitDelay time.Duration

	wg sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

func NewOrderProcessorPool(
	orderRepository repositories.OrderRepository,
	externalAccrualRepository repositories.ExternalAccrualRepository,
	balanceRepository repositories.BalanceRepository,
	logger logger.Logger,
) *OrderProcessorPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &OrderProcessorPool{
		orderRepository:           orderRepository,
		externalAccrualRepository: externalAccrualRepository,
		balanceRepository:         balanceRepository,
		logger:                    logger,
		orderQueue:                make(chan models.OrderModel, 100),
		workerCount:               5,
		pollInterval:              3 * time.Second,
		batchSize:                 20,
		rateLimitDelay:            60 * time.Second,
		ctx:                       ctx,
		cancel:                    cancel,
	}
}

func (p *OrderProcessorPool) Start(ctx context.Context) {
	p.ctx, p.cancel = context.WithCancel(ctx)

	p.logger.Info("Starting order processor pool", "workers", p.workerCount)
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}

	p.wg.Add(1)
	go p.orderLoader()

	p.wg.Wait()
	p.logger.Info("Order processor pool stopped")
}

func (p *OrderProcessorPool) Stop() {
	p.logger.Info("Stopping order processor pool")
	p.cancel()
	close(p.orderQueue)
	p.wg.Wait()
}

func (p *OrderProcessorPool) orderLoader() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	p.loadOrders()

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info("Order loader stopped", "reason", p.ctx.Err())
			return
		case <-ticker.C:
			p.loadOrders()
		}
	}
}

func (p *OrderProcessorPool) loadOrders() {
	orders, err := p.orderRepository.GetPendingOrders(p.ctx, p.batchSize)
	if err != nil {
		p.logger.Error("Failed to get pending orders", "error", err)
		return
	}

	if len(orders) == 0 {
		return
	}

	p.logger.Info("Loaded orders for processing", "count", len(orders))

	for _, order := range orders {
		select {
		case p.orderQueue <- order:
		case <-p.ctx.Done():
			return
		default:
			p.logger.Warn("Order queue is full, skipping order", "order_number", order.OrderNumber)
		}
	}
}

func (p *OrderProcessorPool) worker(id int) {
	defer p.wg.Done()

	p.logger.Info("Worker started", "worker_id", id)

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info("Worker stopped", "worker_id", id, "reason", p.ctx.Err())
			return

		case order, ok := <-p.orderQueue:
			if !ok {
				p.logger.Info("Order queue closed, worker stopping", "worker_id", id)
				return
			}

			if err := p.processOrder(p.ctx, order); err != nil {
				p.logger.Error("Failed to process order",
					"worker_id", id,
					"order_number", order.OrderNumber,
					"error", err,
				)
			} else {
				p.logger.Debug("Order processed successfully",
					"worker_id", id,
					"order_number", order.OrderNumber,
				)
			}
		}
	}
}

func (p *OrderProcessorPool) processOrder(ctx context.Context, order models.OrderModel) error {
	scoring, err := p.externalAccrualRepository.GetScoring(ctx, order.OrderNumber)

	var externalAccrualRepositoryError repositories.ExternalAccrualRepositoryError
	if errors.As(err, &externalAccrualRepositoryError) {
		switch externalAccrualRepositoryError.Code {
		case repositories.ExternalAccrualRepositoryErrorOrderTooManyRequests:
			p.logger.Warn("Rate limit hit for order", "order_number", order.OrderNumber)

			go func() {
				time.Sleep(p.rateLimitDelay)
				select {
				case p.orderQueue <- order:
				case <-p.ctx.Done():
				default:
					p.logger.Warn("Order queue full, order will be retried later", "order_number", order.OrderNumber)
				}
			}()
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

func (p *OrderProcessorPool) accrualBalance(ctx context.Context, order models.OrderModel) error {
	return p.balanceRepository.Accrual(ctx, order.UserID, order.OrderID, order.Accrual)
}

func (p *OrderProcessorPool) GetQueueSize() int {
	return len(p.orderQueue)
}
