package models

type ScoringStatus string

const (
	ScoringStatusRegistered = "REGISTERED"
	ScoringStatusInvalid    = "INVALID"
	ScoringStatusProcessing = "PROCESSING"
	ScoringStatusProcessed  = "PROCESSED"
)

type ScoringModel struct {
	OrderNumber string
	Status      ScoringStatus
	Accrual     float64
}

var EmptyScoringModel = ScoringModel{}
