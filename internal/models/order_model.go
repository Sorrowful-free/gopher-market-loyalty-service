package models

import "time"

type OrderStatus string

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type OrderModel struct {
	Order     string      `json:"order"`
	Status    OrderStatus `json:"status"`
	Accrual   float64     `json:"accrual"`
	CreatedAt time.Time   `json:"uploaded_at"`
}

var EMPTY_ORDER_MODEL = OrderModel{}

func NewOrderModel(order string, status OrderStatus, accrual float64) *OrderModel {
	return &OrderModel{
		Order:     order,
		Status:    status,
		Accrual:   accrual,
		CreatedAt: time.Now(),
	}
}
