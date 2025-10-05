package handlers

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/labstack/echo"
)

type EchoHandlers struct {
	echoRouter     *echo.Echo
	publicGroup    *echo.Group
	protectedGroup *echo.Group

	logger logger.Logger

	authMiddleware *middlewares.EchoAuthMiddleware
}

func NewEchoHandlers(logger logger.Logger) *EchoHandlers {

	echoRouter := echo.New()
	return &EchoHandlers{echoRouter: echoRouter, logger: logger}
}

func (h *EchoHandlers) BuildGroups() {

	h.publicGroup = h.echoRouter.Group(PublicGroup)
	h.protectedGroup = h.echoRouter.Group(ProtectedGroup)
}

func (h *EchoHandlers) BuildAuthMiddleware(jwtSecret string) {

	h.authMiddleware = middlewares.NewEchoAuthMiddleware(jwtSecret, h.logger)

	h.protectedGroup.Use(h.authMiddleware.RequireAuth)
}

func (h *EchoHandlers) BuildRoutes() {

	h.publicGroup.POST(RegisterPath, RegisterHandler)
	h.publicGroup.POST(LoginPath, Login)

	h.protectedGroup.POST(OrdersPath, Orders)
	h.protectedGroup.GET(OrdersPath, OrdersList)
	h.protectedGroup.GET(BalancePath, Balance)
	h.protectedGroup.POST(WithdrawPath, Withdraw)
	h.protectedGroup.GET(WithdrawalsPath, Withdrawals)
}

func (h *EchoHandlers) Run() error {

	return h.echoRouter.Start(":8080")
}
