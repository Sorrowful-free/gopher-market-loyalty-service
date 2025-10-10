package services

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
)

type BalanceService interface {
	GetBalance(userID string) (models.BalanceModel, error)
	Withdraw(userID string, order string, sum float64) error
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

func (s *BalanceServiceImpl) Withdraw(userID string, order string, sum float64) error {
	// s.orderRepository.GetOrder(order)
	return nil
}

func (s *BalanceServiceImpl) GetWithdrawals(userID string) ([]models.WithdrawalModel, error) {
	return nil, nil
}
