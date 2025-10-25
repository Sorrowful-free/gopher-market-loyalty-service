package models

import "time"

type WithdrawalModel struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

var EMPTY_WITHDRAWAL_MODEL = WithdrawalModel{}
var EMPTY_ARRAY_OF_WITHDRAWAL_MODEL = []WithdrawalModel{}

func NewWithdrawalModel(order string, sum float64, processedAt time.Time) *WithdrawalModel {
	return &WithdrawalModel{
		Order:       order,
		Sum:         sum,
		ProcessedAt: processedAt,
	}
}
