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
	OrderID   int         `json:"order"`
	Status    OrderStatus `json:"status"`
	Accrual   float64     `json:"accrual"`
	CreatedAt time.Time   `json:"uploaded_at"`
}

var EMPTY_ORDER_MODEL = OrderModel{}

func NewOrderModel(orderID int, status OrderStatus, accrual float64) *OrderModel {
	return &OrderModel{
		OrderID:   orderID,
		Status:    status,
		Accrual:   accrual,
		CreatedAt: time.Now(),
	}
}
