package handlers

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateOrderHandler(t *testing.T) {
	fiberHandlers := SetupMockFiberHandlers(t)
	fiberApp := fiberHandlers.fiberApp
	orderService := fiberHandlers.orderService

	t.Run("successful_already_created_order", func(t *testing.T) {

		orderService.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.EmptyOrderModel, services.NewOrderServiceError(services.OrderServiceErrorOrderAlreadyExists, "Order already exists"))
		fiberHandlers.MakeValideAuthRequest(t)

		req := httptest.NewRequest(fiber.MethodPost, TestCreateOrderPath, bytes.NewBuffer([]byte(TestValidOrderNumber)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("successful_create_order", func(t *testing.T) {

		orderService.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.EmptyOrderModel, nil)
		fiberHandlers.MakeValideAuthRequest(t)

		req := httptest.NewRequest(fiber.MethodPost, TestCreateOrderPath, bytes.NewBuffer([]byte(TestValidOrderNumber)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusAccepted, resp.StatusCode)
	})

	t.Run("failed_create_order_with_other_user_order", func(t *testing.T) {

		orderService.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.EmptyOrderModel, services.NewOrderServiceError(services.OrderServiceErrorOrderCreatedOtherUser, "Order created other user"))
		fiberHandlers.MakeValideAuthRequest(t)

		req := httptest.NewRequest(fiber.MethodPost, TestCreateOrderPath, bytes.NewBuffer([]byte(TestValidOrderNumber)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})

	t.Run("failed_create_order_with_invalid_order", func(t *testing.T) {

		orderService.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.EmptyOrderModel, services.NewOrderServiceError(services.OrderServiceErrorOrderIDIsInvalid, "Order invalid"))
		fiberHandlers.MakeValideAuthRequest(t)

		req := httptest.NewRequest(fiber.MethodPost, TestCreateOrderPath, bytes.NewBuffer([]byte(TestInvalidOrderNumber)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("failed_create_order_with_internal_error", func(t *testing.T) {

		orderService.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.EmptyOrderModel, errors.New("internal server error"))
		fiberHandlers.MakeValideAuthRequest(t)

		req := httptest.NewRequest(fiber.MethodPost, TestCreateOrderPath, bytes.NewBuffer([]byte(TestValidOrderNumber)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

}
