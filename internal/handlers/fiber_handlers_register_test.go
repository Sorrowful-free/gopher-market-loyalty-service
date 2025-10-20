package handlers

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegisterHandler(t *testing.T) {

	fiberHandlers := SetupMockFiberHandlers(t)
	fiberApp := fiberHandlers.fiberApp
	jwtService := fiberHandlers.jwtService
	userService := fiberHandlers.userService

	t.Run("successful_registration", func(t *testing.T) {

		userService.EXPECT().Register(gomock.Any(), gomock.Any()).Return("userID", nil)
		jwtService.EXPECT().GenerateToken(gomock.Any()).Return("token", nil)

		req := httptest.NewRequest(fiber.MethodPost, TestRegisterUserPath, bytes.NewBuffer([]byte(TestLoginJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}
		defer resp.Body.Close()

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("failed_registration_with_wrong_format", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, TestRegisterUserPath, bytes.NewBuffer([]byte(TestLoginText)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failed_registration_with_existing_login", func(t *testing.T) {

		userService.EXPECT().Register(gomock.Any(), gomock.Any()).Return("userID", services.NewUserServiceError(services.UserServiceErrorUserExists, "User already exists"))

		req := httptest.NewRequest(fiber.MethodPost, TestRegisterUserPath, bytes.NewBuffer([]byte(TestLoginJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})

	t.Run("failed_registration_with_internal_error", func(t *testing.T) {

		userService.EXPECT().Register(gomock.Any(), gomock.Any()).Return("userID", errors.New("internal server error"))

		req := httptest.NewRequest(fiber.MethodPost, TestRegisterUserPath, bytes.NewBuffer([]byte(TestLoginJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

}
