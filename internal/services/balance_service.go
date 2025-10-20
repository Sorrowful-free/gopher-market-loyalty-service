package services

import (
	"errors"

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
	balance, err := s.userRepository.GetBalance(userID)
	if err != nil {
		return models.EMPTY_BALANCE_MODEL, err
	}

	var userRepositoryError repositories.UserRepositoryError
	if errors.As(err, &userRepositoryError) {
		switch userRepositoryError.Code {
		case repositories.UserRepositoryErrorUserNotFound:
			return models.EMPTY_BALANCE_MODEL, NewBalanceServiceError(BalanceServiceErrorUserNotFound, "User not found")
		}
	}
	return balance, nil
}

func (s *BalanceServiceImpl) Withdraw(userID string, orderID string, sum float64) error {

	if !utils.ValidateLuhn(orderID) {
		return NewBalanceServiceError(BalanceServiceErrorOrderIdIsInvalid, "Order id is invalid")
	}

	balance, err := s.userRepository.GetBalance(userID)
	var userRepositoryError repositories.UserRepositoryError
	if errors.As(err, &userRepositoryError) {
		switch userRepositoryError.Code {
		case repositories.UserRepositoryErrorUserNotFound:
			return NewBalanceServiceError(BalanceServiceErrorUserNotFound, "User not found")
		}
	}

	if err != nil {
		return NewBalanceServiceError(BalanceServiceErrorInternalError, "Internal server error")
	}

	if balance.Current < sum {
		return NewBalanceServiceError(BalanceServiceErrorNotEnoughBalance, "Not enough balance")
	}

	return nil
}

func (s *BalanceServiceImpl) GetWithdrawals(userID string) ([]models.WithdrawalModel, error) {
	return nil, nil
}
