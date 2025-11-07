package services

type ExternalAccrualService interface {
	GetScroing(orderID int) (float64, error)
}

type ExternalAccrualServiceImpl struct {
	accrualSystemAddress string
}

func NewExternalAccrualService(accrualSystemAddress string) ExternalAccrualService {
	return &ExternalAccrualServiceImpl{accrualSystemAddress: accrualSystemAddress}
}

func (s *ExternalAccrualServiceImpl) GetScroing(orderID int) (float64, error) {
	return 0, nil
}
