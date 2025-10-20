package models

type BalanceModel struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

var EMPTY_BALANCE_MODEL = BalanceModel{}

func NewBalanceModel(current float64, withdrawn float64) *BalanceModel {
	return &BalanceModel{
		Current:   current,
		Withdrawn: withdrawn,
	}
}
