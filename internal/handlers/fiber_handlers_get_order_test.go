package handlers

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetOrderHandler(t *testing.T) {
	fiberHandlers := SetupMockFiberHandlers(t)
	fiberApp := fiberHandlers.fiberApp
	orderService := fiberHandlers.orderService

	t.Run("successful_get_order", func(t *testing.T) {
		orderService.EXPECT().GetOrder(gomock.Any()).Return(models.OrderModel{
			Order:   TestOrderID,
			Status:  models.OrderStatusNew,
			Accrual: 100,
		}, nil)
		req := httptest.NewRequest(fiber.MethodGet, TestGetOrderPath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("successful_get_order_with_order_not_found", func(t *testing.T) {
		orderService.EXPECT().GetOrder(gomock.Any()).Return(models.OrderModel{}, services.NewOrderServiceError(services.OrderServiceErrorOrderNotFound, "Order not found"))
		req := httptest.NewRequest(fiber.MethodGet, TestGetOrderPath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("failed_get_order_with_too_many_requests", func(t *testing.T) {
		orderService.EXPECT().GetOrder(gomock.Any()).Return(models.OrderModel{}, services.NewOrderServiceError(services.OrderServiceErrorOrderTooManyRequests, "Too many requests"))
		req := httptest.NewRequest(fiber.MethodGet, TestGetOrderPath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})

	t.Run("failed_get_order_with_internal_error", func(t *testing.T) {
		orderService.EXPECT().GetOrder(gomock.Any()).Return(models.OrderModel{}, errors.New("internal server error"))
		req := httptest.NewRequest(fiber.MethodGet, TestGetOrderPath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

}
