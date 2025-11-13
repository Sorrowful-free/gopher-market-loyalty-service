package services

import (
	"context"
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestOrderService(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderRepository := repositories.NewMockOrderRepository(ctrl)
	externalAccrualRepository := repositories.NewMockExternalAccrualRepository(ctrl)
	orderService := NewOrderService(orderRepository, externalAccrualRepository)

	t.Run("successful_create_order", func(t *testing.T) {

		orderRepository.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(models.OrderModel{
			OrderNumber: TestValidOrderID,
			Status:      models.OrderStatusProcessed,
			Accrual:     TestAccrual,
		}, nil)
		externalAccrualRepository.EXPECT().GetScoring(gomock.Any(), gomock.Any()).Return(models.ScoringModel{
			Status:  models.ScoringStatusProcessed,
			Accrual: TestAccrual,
		}, nil)
		orderRepository.EXPECT().UpdateOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(models.OrderModel{
			OrderNumber: TestValidOrderID,
			Status:      models.OrderStatusProcessed,
			Accrual:     TestAccrual,
		}, nil)

		order, err := orderService.CreateOrder(context.TODO(), TestUserID, TestValidOrderID)
		require.NoError(t, err)
		require.Equal(t, TestValidOrderID, order.OrderNumber)
		require.Equal(t, models.OrderStatusProcessed, order.Status)
		require.Equal(t, TestAccrual, order.Accrual)
	})

	t.Run("successful_get_orders_list", func(t *testing.T) {
		externalAccrualRepository.EXPECT().GetScoring(gomock.Any(), gomock.Any()).Return(models.ScoringModel{
			Status:  models.ScoringStatusProcessed,
			Accrual: TestAccrual,
		}, nil)
		orderRepository.EXPECT().GetOrdersList(gomock.Any(), gomock.Any()).Return([]models.OrderModel{
			{
				OrderNumber: TestValidOrderID,
				Status:      models.OrderStatusNew,
				Accrual:     TestAccrual,
			},
		}, nil)
		orderRepository.EXPECT().UpdateOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(models.OrderModel{
			OrderNumber: TestValidOrderID,
			Status:      models.OrderStatusProcessed,
			Accrual:     TestAccrual,
		}, nil)

		orders, err := orderService.GetOrdersList(context.TODO(), TestUserID)
		require.NoError(t, err)
		require.Equal(t, TestValidOrderID, orders[0].OrderNumber)
		require.Equal(t, models.OrderStatusProcessed, orders[0].Status)
		require.Equal(t, TestAccrual, orders[0].Accrual)
	})

	t.Run("successful_get_order_with_order_not_found", func(t *testing.T) {
		orderRepository.EXPECT().GetOrder(gomock.Any(), gomock.Any()).Return(models.EmptyOrderModel, repositories.NewOrderRepositoryError(repositories.OrderRepositoryErrorOrderNotFound, "Order not found"))
		_, err := orderService.GetOrder(context.TODO(), TestValidOrderID)
		var orderServiceError OrderServiceError
		require.ErrorAs(t, err, &orderServiceError)
		require.Equal(t, OrderServiceErrorOrderNotFound, orderServiceError.Code)
		require.Equal(t, "Order not found", orderServiceError.Message)
	})

	t.Run("failed_get_order_with_order_id_is_invalid", func(t *testing.T) {
		_, err := orderService.GetOrder(context.TODO(), TestInvalidOrderID)
		var orderServiceError OrderServiceError
		require.ErrorAs(t, err, &orderServiceError)
		require.Equal(t, OrderServiceErrorOrderIDIsInvalid, orderServiceError.Code)
		require.Equal(t, "Order id is invalid", orderServiceError.Message)
	})

	t.Run("failed_crete_order_is_created_by_other_user", func(t *testing.T) {
		orderRepository.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(models.EmptyOrderModel, repositories.NewOrderRepositoryError(repositories.OrderRepositoryErrorOrderCreatedOtherUser, "Order created by other user"))
		_, err := orderService.CreateOrder(context.TODO(), TestUserID, TestValidOrderID)
		var orderServiceError OrderServiceError
		require.ErrorAs(t, err, &orderServiceError)
		require.Equal(t, OrderServiceErrorOrderCreatedOtherUser, orderServiceError.Code)
		require.Equal(t, "Order created by other user", orderServiceError.Message)
	})
}
