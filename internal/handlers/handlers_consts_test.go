package handlers

const (
	TestRegisterUserPath = "/api/user/register"
	TestLoginPath        = "/api/user/login"
	TestOrderPath        = "/api/user/orders"
	TestBalancePath      = "/api/user/balance"
	TestWithdrawPath     = "/api/user/balance/withdraw"
	TestWithdrawalsPath  = "/api/user/balance/withdrawals"

	TestLoginJSON = `{"login": "test", "password": "test"}`
	TestLoginText = `"login": "test", "password": "test"`

	TestOrderText      = `12345678903`
	TestOrdersListJSON = `[{"number": "12345678903", "status": "NEW", "accrual": 100, "uploaded_at": "2021-01-01T00:00:00Z"}]`

	TestWithdrawJSON = `{"order": "12345678903", "sum": 100}`
)
