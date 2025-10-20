package handlers

import (
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
)

type FiberHandlersMock struct {
	logger         logger.Logger
	jwtService     *services.MockJWTService
	userService    *services.MockUserService
	orderService   *services.MockOrderService
	balanceService *services.MockBalanceService
	fiberApp       *fiber.App
}

func SetupMockFiberHandlers(t *testing.T) *FiberHandlersMock {

	logger := logger.NewZapLogger()
	ctrl := gomock.NewController(t)
	jwtService := services.NewMockJWTService(ctrl)
	userService := services.NewMockUserService(ctrl)
	orderService := services.NewMockOrderService(ctrl)
	balanceService := services.NewMockBalanceService(ctrl)

	fiberHandlers := NewFiberHandlers(logger, jwtService, userService, orderService, balanceService)

	fiberHandlers.BuildGroups()
	fiberHandlers.BuildAuthMiddleware()
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
