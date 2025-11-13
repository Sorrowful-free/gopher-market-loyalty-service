package models

type WithdrawRequest struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
}
