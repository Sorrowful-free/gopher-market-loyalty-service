package models

type GetOrderResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int64  `json:"accrual"`
}
