package models

type GetBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn int32   `json:"withdrawn"`
}

func NewGetBalanceResponse(current float64, withdrawn int32) *GetBalanceResponse {
	return &GetBalanceResponse{
		Current:   current,
		Withdrawn: withdrawn,
	}
}
