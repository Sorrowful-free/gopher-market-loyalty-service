package services

import (
	"context"
	"errors"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/utils"
)

//go:generate mockgen -source=balance_service.go -destination=mock_balance_service.go -package=services
type BalanceService interface {
	GetBalance(ctx context.Context, userID string) (models.BalanceModel, error)
	Withdraw(ctx context.Context, userID string, orderID string, sum float64) error
	GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawalModel, error)
}

type BalanceServiceImpl struct {
	userRepository  repositories.UserRepository
	orderRepository repositories.OrderRepository
}

func NewBalanceService(userRepository repositories.UserRepository, orderRepository repositories.OrderRepository) BalanceService {
	return &BalanceServiceImpl{userRepository: userRepository, orderRepository: orderRepository}
}

func (s *BalanceServiceImpl) GetBalance(ctx context.Context, userID string) (models.BalanceModel, error) {
	if ctx.Err() != nil {
		return models.EMPTY_BALANCE_MODEL, ctx.Err()
	}

	balance, err := s.userRepository.GetBalance(ctx, userID)
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

func (s *BalanceServiceImpl) Withdraw(ctx context.Context, userID string, orderID string, sum float64) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if !utils.ValidateLuhn(orderID) {
		return NewBalanceServiceError(BalanceServiceErrorOrderIdIsInvalid, "Order id is invalid")
	}

	balance, err := s.userRepository.GetBalance(ctx, userID)
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

func (s *BalanceServiceImpl) GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawalModel, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// withdrawals, err := s.orderRepository.GetWithdrawals(ctx, userID)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
