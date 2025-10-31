package models

type WithdrawRequest struct {
	Order int     `json:"order"`
	Sum   float64 `json:"sum"`
}
