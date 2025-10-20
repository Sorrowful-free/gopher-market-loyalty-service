package services

import (
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestOrderService(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderRepository := repositories.NewMockOrderRepository(ctrl)
	orderService := NewOrderService(orderRepository)

	t.Run("successful_create_order", func(t *testing.T) {
		orderRepository.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(models.OrderModel{
			Order:   TestValidOrderID,
			Status:  models.OrderStatusNew,
			Accrual: 100,
		}, nil)
		order, err := orderService.CreateOrder(TestUserID, TestValidOrderID)
		require.NoError(t, err)
		require.Equal(t, TestValidOrderID, order.Order)
	})

	t.Run("successful_get_orders_list", func(t *testing.T) {
		orderRepository.EXPECT().GetOrdersList(gomock.Any()).Return([]models.OrderModel{
			{
				Order:   TestValidOrderID,
				Status:  models.OrderStatusNew,
				Accrual: 100,
			},
		}, nil)

		orders, err := orderService.GetOrdersList(TestUserID)
		require.NoError(t, err)
		require.Equal(t, TestValidOrderID, orders[0].Order)
	})

	t.Run("successful_get_order", func(t *testing.T) {
		orderRepository.EXPECT().GetOrder(gomock.Any()).Return(models.OrderModel{
			Order:   TestValidOrderID,
			Status:  models.OrderStatusNew,
			Accrual: 100,
		}, nil)

		order, err := orderService.GetOrder(TestValidOrderID)
		require.NoError(t, err)
		require.Equal(t, TestValidOrderID, order.Order)
	})

	t.Run("successful_get_order_with_order_not_found", func(t *testing.T) {
		orderRepository.EXPECT().GetOrder(gomock.Any()).Return(models.EMPTY_ORDER_MODEL, repositories.NewOrderRepositoryError(repositories.OrderRepositoryErrorOrderNotFound, "Order not found"))
		_, err := orderService.GetOrder(TestValidOrderID)
		var orderServiceError OrderServiceError
		require.ErrorAs(t, err, &orderServiceError)
		require.Equal(t, OrderServiceErrorOrderNotFound, orderServiceError.Code)
		require.Equal(t, "Order not found", orderServiceError.Message)
	})

	t.Run("failed_get_order_with_order_id_is_invalid", func(t *testing.T) {
		_, err := orderService.GetOrder(TestInvalidOrderID)
		var orderServiceError OrderServiceError
		require.ErrorAs(t, err, &orderServiceError)
		require.Equal(t, OrderServiceErrorOrderIdIsInvalid, orderServiceError.Code)
		require.Equal(t, "Order id is invalid", orderServiceError.Message)
	})

	t.Run("failed_crete_order_is_created_by_other_user", func(t *testing.T) {
		orderRepository.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(models.EMPTY_ORDER_MODEL, repositories.NewOrderRepositoryError(repositories.OrderRepositoryErrorOrderCreatedOtherUser, "Order created by other user"))
		_, err := orderService.CreateOrder(TestUserID, TestValidOrderID)
		var orderServiceError OrderServiceError
		require.ErrorAs(t, err, &orderServiceError)
		require.Equal(t, OrderServiceErrorOrderCreatedOtherUser, orderServiceError.Code)
		require.Equal(t, "Order created by other user", orderServiceError.Message)
	})
}
