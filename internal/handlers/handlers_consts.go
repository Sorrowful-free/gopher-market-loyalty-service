package handlers

const (
	UserGroup = "/api/user"

	RegisterUserPath = "/register" //POST api/user/register
	LoginUserPath    = "/login"    //POST api/user/login

	OrderGroup = "/orders" //api/user/orders

	CreateOrderPath   = "/"       //POST api/user/orders
	GetOrdersListPath = "/"       //GET api/user/orders
	GetOrderPath      = "/:order" //GET api/user/orders/:order

	BalanceGroup = "/balance" //api/user/balance

	GetBalancePath  = "/"            //GET api/user/balance
	WithdrawPath    = "/withdraw"    //POST api/user/balance/withdraw
	WithdrawalsPath = "/withdrawals" //GET api/user/balance/withdrawals
)
