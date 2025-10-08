package handlers

const (
	PublicGroup    = "/api/user"
	ProtectedGroup = "/api/user"

	RegisterPath      = "/api/user/register"
	LoginPath         = "/api/user/login"
	CreateOrderPath   = "/api/user/orders"
	GetOrdersListPath = "/api/user/orders"
	GetOrderPath      = "/api/user/orders/:order"

	BalancePath  = "/api/user/balance"
	WithdrawPath = "/api/user/balance/withdraw"

	WithdrawalsPath = "/api/user/balance/withdrawals"
)
