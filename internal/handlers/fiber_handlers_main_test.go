package handlers

import (
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
)

func TestSetupFiberHandlers(t *testing.T) *FiberHandlers {
	var logger logger.Logger
	var jwtService services.JWTService
	var userService services.UserService
	var orderService services.OrderService

	return NewFiberHandlers(logger, jwtService, userService, orderService)
}
