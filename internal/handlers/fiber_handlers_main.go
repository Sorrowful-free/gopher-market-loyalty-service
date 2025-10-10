package handlers

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/skip"
)

type FiberHandlers struct {
	fiberApp     *fiber.App
	userGroup    fiber.Router
	orderGroup   fiber.Router
	balanceGroup fiber.Router

	logger       logger.Logger
	jwtService   services.JWTService
	userService  services.UserService
	orderService services.OrderService

	authMiddleware *middlewares.FiberAuthMiddleware
}

func NewFiberHandlers(logger logger.Logger, jwtService services.JWTService, userService services.UserService, orderService services.OrderService) *FiberHandlers {
	fiberApp := fiber.New()
	return &FiberHandlers{fiberApp: fiberApp, logger: logger, jwtService: jwtService, userService: userService, orderService: orderService}
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

	h.orderGroup.Use(skip.New(h.authMiddleware.RequireAuth, func(c *fiber.Ctx) bool {
		return c.Path() == GetOrderPath
	}))
	h.orderGroup.Post(CreateOrderPath, middlewares.ValidateRequestAsText(), h.CreateOrderHandler)
	h.orderGroup.Get(GetOrdersListPath, GetOrdersList)
	h.orderGroup.Get(GetOrderPath, GetOrder)

	h.balanceGroup.Use(h.authMiddleware.RequireAuth)
	h.balanceGroup.Get(GetBalancePath, Balance)
	h.balanceGroup.Post(WithdrawBalancePath, Withdraw)
	h.balanceGroup.Get(BalanceWithdrawalsPath, Withdrawals)
}

func (h *FiberHandlers) Run() error {
	return h.fiberApp.Listen(":8080")
}
