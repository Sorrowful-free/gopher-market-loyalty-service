package handlers

const (
	UserGroup = "/api/user"

	RegisterUserPath = "/register"
	LoginUserPath    = "/login"

	OrderGroup = "/order"

	CreateOrderPath   = "/"
	GetOrdersListPath = "/"
	GetOrderPath      = "/:order"

	BalanceGroup = "/balance"

	GetBalancePath         = "/"
	WithdrawBalancePath    = "/withdraw"
	BalanceWithdrawalsPath = "/withdrawals"
)
