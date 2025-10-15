package services

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/utils"
)

type BalanceService interface {
	GetBalance(userID string) (models.BalanceModel, error)
	Withdraw(userID string, orderID string, sum float64) error
	GetWithdrawals(userID string) ([]models.WithdrawalModel, error)
}

type BalanceServiceImpl struct {
	userRepository  repositories.UserRepository
	orderRepository repositories.OrderRepository
}

func NewBalanceService(userRepository repositories.UserRepository, orderRepository repositories.OrderRepository) BalanceService {
	return &BalanceServiceImpl{userRepository: userRepository, orderRepository: orderRepository}
}

func (s *BalanceServiceImpl) GetBalance(userID string) (models.BalanceModel, error) {
	return s.userRepository.GetBalance(userID)
}

func (s *BalanceServiceImpl) Withdraw(userID string, orderID string, sum float64) error {
	if !utils.ValidateLuhn(orderID) {
		return NewBalanceServiceError(BalanceServiceErrorOrderIdIsInvalid, "Order id is invalid")
	}

	// s.orderRepository.GetOrder(order)
	return nil
}

func (s *BalanceServiceImpl) GetWithdrawals(userID string) ([]models.WithdrawalModel, error) {
	return nil, nil
}
