package models

import "time"

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type OrderModel struct {
	Order     string    `json:"order"`
	Status    string    `json:"status"`
	Accrual   int64     `json:"accrual"`
	CreatedAt time.Time `json:"uploaded_at"`
}

func NewOrderModel(order string, status string, accrual int64) *OrderModel {
	return &OrderModel{
		Order:     order,
		Status:    status,
		Accrual:   accrual,
		CreatedAt: time.Now(),
	}
}
