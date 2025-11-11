package models

type ScoringStatus string

const (
	ScoringStatusRegistered ScoringStatus = "REGISTERED"
	ScoringStatusInvalid    ScoringStatus = "INVALID"
	ScoringStatusProcessing ScoringStatus = "PROCESSING"
	ScoringStatusProcessed  ScoringStatus = "PROCESSED"
)

type ScoringModel struct {
	OrderNumber string
	Status      ScoringStatus
	Accrual     float64
}

var EmptyScoringModel = ScoringModel{}
