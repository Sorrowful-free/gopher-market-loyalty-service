package models

import "time"

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type OrderModel struct {
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	Accrual   int64     `json:"accrual"`
	CreatedAt time.Time `json:"uploaded_at"`
}

func NewOrderModel(number string, status string, accrual int64) *OrderModel {
	return &OrderModel{
		Number:    number,
		Status:    status,
		Accrual:   accrual,
		CreatedAt: time.Now(),
	}
}
