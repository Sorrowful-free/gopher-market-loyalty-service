package models

type GetOrderResponse struct {
	OrderNumber string      `json:"order"`
	Status      OrderStatus `json:"status"`
	Accrual     float64     `json:"accrual"`
}
