package handlers

import (
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
)

type FiberHandlersMock struct {
	logger         logger.Logger
	jwtService     *mocks.MockJWTService
	userService    *mocks.MockUserService
	orderService   *mocks.MockOrderService
	balanceService *mocks.MockBalanceService
	fiberApp       *fiber.App
}

func SetupMockFiberHandlers(t *testing.T) *FiberHandlersMock {

	logger := logger.NewZapLogger()
	ctrl := gomock.NewController(t)
	jwtService := mocks.NewMockJWTService(ctrl)
	userService := mocks.NewMockUserService(ctrl)
	orderService := mocks.NewMockOrderService(ctrl)
	balanceService := mocks.NewMockBalanceService(ctrl)

	fiberHandlers := NewFiberHandlers(logger, jwtService, userService, orderService, balanceService)

	fiberHandlers.BuildGroups()
	fiberHandlers.BuildAuthMiddleware("test")
	fiberHandlers.BuildRoutes()

	return &FiberHandlersMock{
		logger:         logger,
		jwtService:     jwtService,
		userService:    userService,
		orderService:   orderService,
		balanceService: balanceService,
		fiberApp:       fiberHandlers.fiberApp,
	}
}
