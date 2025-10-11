package handlers

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

type FiberHandlers struct {
	fiberApp     *fiber.App
	userGroup    fiber.Router
	orderGroup   fiber.Router
	balanceGroup fiber.Router

	logger         logger.Logger
	jwtService     services.JWTService
	userService    services.UserService
	orderService   services.OrderService
	balanceService services.BalanceService

	authMiddleware *middlewares.FiberAuthMiddleware
}

func NewFiberHandlers(logger logger.Logger, jwtService services.JWTService, userService services.UserService, orderService services.OrderService, balanceService services.BalanceService) *FiberHandlers {
	fiberApp := fiber.New()
	return &FiberHandlers{fiberApp: fiberApp, logger: logger, jwtService: jwtService, userService: userService, orderService: orderService, balanceService: balanceService}
}

func (h *FiberHandlers) BuildGroups() {
	h.userGroup = h.fiberApp.Group(UserGroup)
	h.orderGroup = h.userGroup.Group(OrderGroup)
	h.balanceGroup = h.userGroup.Group(BalanceGroup)
}

func (h *FiberHandlers) BuildAuthMiddleware(jwtSecret string) {
	h.authMiddleware = middlewares.NewFiberAuthMiddleware(jwtSecret, h.logger, h.jwtService)
}

func (h *FiberHandlers) BuildRoutes() {

	h.userGroup.Post(RegisterUserPath, middlewares.ValidateRequestAsJSON(models.RegisterRequest{}), h.RegisterHandler)
	h.userGroup.Post(LoginUserPath, middlewares.ValidateRequestAsJSON(models.LoginRequest{}), h.LoginHandler)

	h.orderGroup.Post(CreateOrderPath, h.authMiddleware.RequireAuth, middlewares.ValidateRequestAsText(), h.CreateOrderHandler)
	h.orderGroup.Get(GetOrdersListPath, h.authMiddleware.RequireAuth, h.GetOrdersListHandler)
	h.orderGroup.Get(GetOrderPath, h.GetOrderHandler)

	h.balanceGroup.Use(h.authMiddleware.RequireAuth)
	h.balanceGroup.Get(GetBalancePath, h.GetBalanceHandler)
	h.balanceGroup.Post(WithdrawPath, middlewares.ValidateRequestAsJSON(models.WithdrawRequest{}), h.WithdrawHandler)
	h.balanceGroup.Get(WithdrawalsPath, h.WithdrawalsHandler)
}

func (h *FiberHandlers) Run() error {
	return h.fiberApp.Listen(":8080")
}
