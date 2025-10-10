package handlers

const (
	TestRegisterUserPath = "/api/user/register"
	TestLoginUserPath    = "/api/user/login"

	TestCreateOrderPath = "/api/user/orders"

	TestLoginJSON = `{"login": "test", "password": "test"}`
	TestLoginText = `"login": "test", "password": "test"`

	TestOrderText      = `12345678903`
	TestOrdersListJSON = `[{"number": "12345678903", "status": "NEW", "accrual": 100, "uploaded_at": "2021-01-01T00:00:00Z"}]`

	TestWithdrawJSON = `{"order": "12345678903", "sum": 100}`
)
