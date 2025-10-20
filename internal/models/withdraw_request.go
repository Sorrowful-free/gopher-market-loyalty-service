package models

type WithdrawRequest struct {
	Order string `json:"order"`
	Sum   int32  `json:"sum"`
}
