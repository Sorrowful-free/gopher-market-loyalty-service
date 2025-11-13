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

// OrderProcessorPool обрабатывает заказы с использованием Worker Pool и очереди задач
type OrderProcessorPool struct {
	orderRepository           repositories.OrderRepository
	externalAccrualRepository repositories.ExternalAccrualRepository
	balanceRepository         repositories.BalanceRepository
	logger                    logger.Logger

	// Очередь задач для обработки заказов
	orderQueue chan models.OrderModel

	// Количество воркеров для параллельной обработки
	workerCount int

	// Интервал между проверками заказов в БД
	pollInterval time.Duration

	// Количество заказов для загрузки за один цикл
	batchSize int

	// Задержка при получении 429 (rate limit)
	rateLimitDelay time.Duration

	// WaitGroup для синхронизации воркеров
	wg sync.WaitGroup

	// Контекст для остановки
	ctx    context.Context
	cancel context.CancelFunc
}

// NewOrderProcessorPool создает новый процессор заказов с Worker Pool
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
		orderQueue:                make(chan models.OrderModel, 100), // Буферизованная очередь на 100 заказов
		workerCount:               5,                                 // 5 воркеров для параллельной обработки
		pollInterval:              3 * time.Second,                   // Проверка каждые 3 секунды
		batchSize:                 20,                                // Загружаем до 20 заказов за раз
		rateLimitDelay:            60 * time.Second,                  // Задержка при rate limit
		ctx:                       ctx,
		cancel:                    cancel,
	}
}

// Start запускает процессор с Worker Pool
func (p *OrderProcessorPool) Start(ctx context.Context) {
	// Объединяем контексты
	p.ctx, p.cancel = context.WithCancel(ctx)

	// Запускаем воркеры
	p.logger.Info("Starting order processor pool", "workers", p.workerCount)
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}

	// Запускаем загрузчик заказов
	p.wg.Add(1)
	go p.orderLoader()

	// Ждем завершения всех горутин
	p.wg.Wait()
	p.logger.Info("Order processor pool stopped")
}

// Stop останавливает процессор
func (p *OrderProcessorPool) Stop() {
	p.logger.Info("Stopping order processor pool")
	p.cancel()
	close(p.orderQueue)
	p.wg.Wait()
}

// orderLoader загружает заказы из БД и помещает их в очередь
func (p *OrderProcessorPool) orderLoader() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	// Загружаем сразу при старте
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

// loadOrders загружает заказы из БД и добавляет их в очередь
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

	// Добавляем заказы в очередь (неблокирующе)
	for _, order := range orders {
		select {
		case p.orderQueue <- order:
			// Заказ добавлен в очередь
		case <-p.ctx.Done():
			return
		default:
			// Очередь переполнена, пропускаем этот заказ
			// Он будет обработан в следующем цикле
			p.logger.Warn("Order queue is full, skipping order", "order_number", order.OrderNumber)
		}
	}
}

// worker обрабатывает заказы из очереди
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
				// Очередь закрыта
				p.logger.Info("Order queue closed, worker stopping", "worker_id", id)
				return
			}

			// Обрабатываем заказ
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

// processOrder обрабатывает один заказ
func (p *OrderProcessorPool) processOrder(ctx context.Context, order models.OrderModel) error {
	// Получаем информацию о статусе заказа из внешней системы
	scoring, err := p.externalAccrualRepository.GetScoring(ctx, order.OrderNumber)

	var externalAccrualRepositoryError repositories.ExternalAccrualRepositoryError
	if errors.As(err, &externalAccrualRepositoryError) {
		switch externalAccrualRepositoryError.Code {
		case repositories.ExternalAccrualRepositoryErrorOrderTooManyRequests:
			// Rate limit - возвращаем заказ в очередь через некоторое время
			p.logger.Warn("Rate limit hit for order", "order_number", order.OrderNumber)

			// Запускаем горутину для повторной попытки через задержку
			go func() {
				time.Sleep(p.rateLimitDelay)
				select {
				case p.orderQueue <- order:
					// Заказ возвращен в очередь
				case <-p.ctx.Done():
					// Процессор остановлен
				default:
					// Очередь переполнена, заказ будет обработан в следующем цикле загрузки
					p.logger.Warn("Order queue full, order will be retried later", "order_number", order.OrderNumber)
				}
			}()
			return nil

		case repositories.ExternalAccrualRepositoryErrorOrderNotRegistered:
			// Заказ еще не зарегистрирован в системе начислений
			// Это нормально, просто пропускаем - он будет обработан позже
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
func (p *OrderProcessorPool) accrualBalance(ctx context.Context, order models.OrderModel) error {
	return p.balanceRepository.Accrual(ctx, order.UserID, order.OrderID, order.Accrual)
}

// GetQueueSize возвращает текущий размер очереди (для мониторинга)
func (p *OrderProcessorPool) GetQueueSize() int {
	return len(p.orderQueue)
}
