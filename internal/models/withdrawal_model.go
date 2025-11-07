package models

import "time"

type WithdrawalModel struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	OrderID     int       `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

var EmptyWithdrawalModel = WithdrawalModel{}
var EmptyArrayOfWithdrawalModel = []WithdrawalModel{}

func NewWithdrawalModel(id int, userID int, orderID int, sum float64, processedAt time.Time) *WithdrawalModel {
	return &WithdrawalModel{
		ID:          id,
		UserID:      userID,
		OrderID:     orderID,
		Sum:         sum,
		ProcessedAt: processedAt,
	}
}
