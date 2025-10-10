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
	CreatedAt time.Time `json:"uploaded_at"`
}

func NewOrderModel(number string, status string) *OrderModel {
	return &OrderModel{
		Number:    number,
		Status:    status,
		CreatedAt: time.Now(),
	}
}
