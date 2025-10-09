package handlers

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

type FiberHandlers struct {
	fiberApp       *fiber.App
	publicGroup    fiber.Router
	protectedGroup fiber.Router

	logger       logger.Logger
	userService  services.UserService
	orderService services.OrderService

	authMiddleware *middlewares.FiberAuthMiddleware
}

func NewFiberHandlers(logger logger.Logger, userService services.UserService, orderService services.OrderService) *FiberHandlers {
	fiberApp := fiber.New()
	return &FiberHandlers{fiberApp: fiberApp, logger: logger, userService: userService, orderService: orderService}
}

func (h *FiberHandlers) BuildGroups() {
	h.publicGroup = h.fiberApp.Group(PublicGroup)
	h.protectedGroup = h.fiberApp.Group(ProtectedGroup)
}

func (h *FiberHandlers) BuildAuthMiddleware(jwtSecret string) {
	h.authMiddleware = middlewares.NewFiberAuthMiddleware(jwtSecret, h.logger)
	h.protectedGroup.Use(h.authMiddleware.RequireAuth)
}

func (h *FiberHandlers) BuildRoutes() {
	h.publicGroup.Post(RegisterPath, RegisterHandler)
	h.publicGroup.Post(LoginPath, Login)
	h.publicGroup.Get(GetOrderPath, GetOrder)

	h.protectedGroup.Post(CreateOrderPath, CreateOrder)
	h.protectedGroup.Get(GetOrdersListPath, GetOrdersList)
	h.protectedGroup.Get(BalancePath, Balance)
	h.protectedGroup.Post(WithdrawPath, Withdraw)
	h.protectedGroup.Get(WithdrawalsPath, Withdrawals)
}

func (h *FiberHandlers) Run() error {
	return h.fiberApp.Listen(":8080")
}
