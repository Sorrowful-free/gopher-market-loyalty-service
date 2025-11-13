package models

type ScoringStatus string

const (
	ScoringStatusRegistered ScoringStatus = "REGISTERED"
	ScoringStatusInvalid    ScoringStatus = "INVALID"
	ScoringStatusProcessing ScoringStatus = "PROCESSING"
	ScoringStatusProcessed  ScoringStatus = "PROCESSED"
)

type ScoringModel struct {
	OrderNumber string        `json:"order"`
	Status      ScoringStatus `json:"status"`
	Accrual     float64       `json:"accrual"`
}

var EmptyScoringModel = ScoringModel{}
