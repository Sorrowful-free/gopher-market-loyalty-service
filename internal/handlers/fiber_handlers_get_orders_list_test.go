package handlers

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetOrdersListHandler(t *testing.T) {
	fiberHandlers := SetupMockFiberHandlers(t)
	fiberApp := fiberHandlers.fiberApp
	orderService := fiberHandlers.orderService
	jwtService := fiberHandlers.jwtService

	t.Run("successful_get_orders_list", func(t *testing.T) {

		orders := []models.OrderModel{
			{
				Number:    TestOrderID,
				Status:    models.OrderStatusNew,
				CreatedAt: time.Now(),
			},
		}

		orderService.EXPECT().GetOrdersList(gomock.Any()).Return(orders, nil)
		jwtService.EXPECT().ValidateToken(gomock.Any()).Return(&services.JWTClaims{}, nil)
		jwtService.EXPECT().ExtractToken(gomock.Any()).Return("userID", nil)

		req := httptest.NewRequest(fiber.MethodGet, TestGetOrdersListPath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("successful_get_orders_list_with_empty_list", func(t *testing.T) {

		orderService.EXPECT().GetOrdersList(gomock.Any()).Return([]models.OrderModel{}, nil)
		jwtService.EXPECT().ValidateToken(gomock.Any()).Return(&services.JWTClaims{}, nil)
		jwtService.EXPECT().ExtractToken(gomock.Any()).Return("userID", nil)

		req := httptest.NewRequest(fiber.MethodGet, TestGetOrdersListPath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("failed_get_orders_list_with_internal_error", func(t *testing.T) {

		orderService.EXPECT().GetOrdersList(gomock.Any()).Return([]models.OrderModel{}, errors.New("internal server error"))
		jwtService.EXPECT().ValidateToken(gomock.Any()).Return(&services.JWTClaims{}, nil)
		jwtService.EXPECT().ExtractToken(gomock.Any()).Return("userID", nil)

		req := httptest.NewRequest(fiber.MethodGet, TestGetOrdersListPath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
