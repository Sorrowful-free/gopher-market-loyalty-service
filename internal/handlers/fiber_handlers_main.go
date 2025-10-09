package handlers

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
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
	jsonValidator  *middlewares.ValidateContentTypeMiddleware
	textValidator  *middlewares.ValidateContentTypeMiddleware
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
	h.jsonValidator = middlewares.NewValidateContentTypeMiddleware(fiber.MIMEApplicationJSON)
	h.textValidator = middlewares.NewValidateContentTypeMiddleware(fiber.MIMETextPlain)
}

func (h *FiberHandlers) BuildRoutes() {

	h.userGroup.Use(h.jsonValidator.ValidateContentType)
	h.userGroup.Post(RegisterUserPath, h.RegisterHandler)
	h.userGroup.Post(LoginUserPath, Login)

	h.orderGroup.Post(CreateOrderPath, CreateOrder).Use(h.textValidator.ValidateContentType, h.authMiddleware.RequireAuth)
	h.orderGroup.Get(GetOrdersListPath, GetOrdersList).Use(h.jsonValidator.ValidateContentType, h.authMiddleware.RequireAuth)
	h.orderGroup.Get(GetOrderPath, GetOrder).Use(h.jsonValidator.ValidateContentType)

	h.balanceGroup.Use(h.authMiddleware.RequireAuth)
	h.balanceGroup.Use(h.jsonValidator.ValidateContentType)
	h.balanceGroup.Get(GetBalancePath, Balance)
	h.balanceGroup.Post(WithdrawBalancePath, Withdraw)
	h.balanceGroup.Get(BalanceWithdrawalsPath, Withdrawals)
}

func (h *FiberHandlers) Run() error {
	return h.fiberApp.Listen(":8080")
}
