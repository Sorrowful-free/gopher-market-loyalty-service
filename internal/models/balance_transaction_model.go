package models

import "time"

type TransactionType string

const (
	TransactionTypeAccrual    TransactionType = "ACCRUAL"
	TransactionTypeWithdrawal TransactionType = "WITHDRAWAL"
)

type BalanceTransactionModel struct {
	ID              int            `json:"id"`
	UserID          int            `json:"user_id"`
	OrderID         *int           `json:"order_id,omitempty"` // NULL для withdrawals
	TransactionType TransactionType `json:"transaction_type"`
	Amount          float64        `json:"amount"`
	CreatedAt       time.Time      `json:"created_at"`
}

var EmptyBalanceTransactionModel = BalanceTransactionModel{}
var EmptyArrayOfBalanceTransactionModel = []BalanceTransactionModel{}

func NewBalanceTransactionModel(
	id int,
	userID int,
	orderID *int,
	transactionType TransactionType,
	amount float64,
	createdAt time.Time,
) *BalanceTransactionModel {
	return &BalanceTransactionModel{
		ID:              id,
		UserID:          userID,
		OrderID:         orderID,
		TransactionType: transactionType,
		Amount:          amount,
		CreatedAt:       createdAt,
	}
}

