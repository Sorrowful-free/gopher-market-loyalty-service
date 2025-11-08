package models

import "time"

type OrderStatus string

const (
	OrderStatusRegistered = "REGISTERED"
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type OrderModel struct {
	OrderID     int         `json:"order_id"`
	OrderNumber string      `json:"number"`
	UserID      int         `json:"user_id"`
	Status      OrderStatus `json:"status"`
	Accrual     float64     `json:"accrual"`
	UploadedAt  time.Time   `json:"uploaded_at"`
}

var EmptyOrderModel = OrderModel{}
var EmptyArrayOfOrderModel = []OrderModel{}

func NewOrderModel(orderID int, userID int, status OrderStatus, accrual float64) *OrderModel {
	return &OrderModel{
		OrderID:    orderID,
		UserID:     userID,
		Status:     status,
		Accrual:    accrual,
		UploadedAt: time.Now().UTC(),
	}
}
