package handlers

const (
	TestRegisterUserPath = "/api/user/register"
	TestLoginUserPath    = "/api/user/login"

	TestCreateOrderPath   = "/api/user/orders"
	TestGetOrdersListPath = "/api/user/orders"
	TestGetOrderPath      = "/api/user/orders/12345678903"

	TestGetBalancePath  = "/api/user/balance"
	TestWithdrawPath    = "/api/user/balance/withdraw"
	TestWithdrawalsPath = "/api/user/balance/withdrawals"

	TestLoginJSON = `{"login": "test", "password": "test"}`
	TestLoginText = `"login": "test", "password": "test"`

	TestUserID             = 1
	TestValidOrderID       = 3
	TestValidOrderNumber   = "12345678903"
	TestInvalidOrderNumber = "12345678904"
	TestWithdrawSum        = float64(100)

	TestWithdrawJSON = `{"order": "12345678903", "sum": 100}`
)
