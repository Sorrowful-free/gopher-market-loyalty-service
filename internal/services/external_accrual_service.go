package services

type ExternalAccrualService interface {
	GetScroing(order string) (int, error)
}

type ExternalAccrualServiceImpl struct {
	accrualSystemAddress string
}

func NewExternalAccrualService(accrualSystemAddress string) ExternalAccrualService {
	return &ExternalAccrualServiceImpl{accrualSystemAddress: accrualSystemAddress}
}

func (s *ExternalAccrualServiceImpl) GetScroing(order string) (int, error) {
	return 0, nil
}
