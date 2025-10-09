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
	h.publicGroup = h.fiberApp.Group(PublicGroup)
	h.protectedGroup = h.fiberApp.Group(ProtectedGroup)
}

func (h *FiberHandlers) BuildAuthMiddleware(jwtSecret string) {
	h.authMiddleware = middlewares.NewFiberAuthMiddleware(jwtSecret, h.logger, h.jwtService)
	h.jsonValidator = middlewares.NewValidateContentTypeMiddleware(fiber.MIMEApplicationJSON)
	h.textValidator = middlewares.NewValidateContentTypeMiddleware(fiber.MIMETextPlain)
}

func (h *FiberHandlers) BuildRoutes() {
	h.publicGroup.Post(RegisterPath, h.RegisterHandler).Use(h.jsonValidator.ValidateContentType)
	h.publicGroup.Post(LoginPath, Login).Use(h.jsonValidator.ValidateContentType)
	h.publicGroup.Get(GetOrderPath, GetOrder).Use(h.jsonValidator.ValidateContentType)

	h.protectedGroup.Use(h.authMiddleware.RequireAuth)
	h.protectedGroup.Post(CreateOrderPath, CreateOrder).Use(h.textValidator.ValidateContentType)
	h.protectedGroup.Get(GetOrdersListPath, GetOrdersList).Use(h.jsonValidator.ValidateContentType)
	h.protectedGroup.Get(BalancePath, Balance).Use(h.jsonValidator.ValidateContentType)
	h.protectedGroup.Post(WithdrawPath, Withdraw).Use(h.jsonValidator.ValidateContentType)
	h.protectedGroup.Get(WithdrawalsPath, Withdrawals).Use(h.jsonValidator.ValidateContentType)
}

func (h *FiberHandlers) Run() error {
	return h.fiberApp.Listen(":8080")
}
