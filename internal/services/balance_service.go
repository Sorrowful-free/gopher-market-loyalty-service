package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/utils"
)

//go:generate mockgen -source=balance_service.go -destination=mock_balance_service.go -package=services
type BalanceService interface {
	GetBalance(ctx context.Context, userID int) (models.BalanceModel, error)
	Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error
	GetWithdrawals(ctx context.Context, userID int) ([]models.WithdrawalModel, error)
}

type BalanceServiceImpl struct {
	balanceRepository repositories.BalanceRepository
}

func NewBalanceService(balanceRepository repositories.BalanceRepository) BalanceService {
	return &BalanceServiceImpl{balanceRepository: balanceRepository}
}

func (s *BalanceServiceImpl) GetBalance(ctx context.Context, userID int) (models.BalanceModel, error) {
	if ctx.Err() != nil {
		return models.EmptyBalanceModel, ctx.Err()
	}

	balance, err := s.balanceRepository.GetBalance(ctx, userID)
	if err != nil {
		return models.EmptyBalanceModel, err
	}

	var userRepositoryError repositories.UserRepositoryError
	if errors.As(err, &userRepositoryError) {
		switch userRepositoryError.Code {
		case repositories.UserRepositoryErrorUserNotFound:
			return models.EmptyBalanceModel, NewBalanceServiceError(BalanceServiceErrorUserNotFound, "User not found")
		}
	}
	return balance, nil
}

func (s *BalanceServiceImpl) Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if !utils.ValidateLuhn(orderNumber) {
		return NewBalanceServiceError(BalanceServiceErrorOrderIDIsInvalid, "Order id is invalid")
	}

	balance, err := s.balanceRepository.GetBalance(ctx, userID)
	var balanceRepositoryError repositories.BalanceRepositoryError
	if errors.As(err, &balanceRepositoryError) {
		switch balanceRepositoryError.Code {
		case repositories.BalanceRepositoryErrorUserNotFound:
			return NewBalanceServiceError(BalanceServiceErrorUserNotFound, "User not found")
		}
	}

	if err != nil {
		return NewBalanceServiceError(BalanceServiceErrorInternalError, "Internal server error")
	}

	if balance.Current < sum {
		return NewBalanceServiceError(BalanceServiceErrorNotEnoughBalance, "Not enough balance")
	}

	return s.balanceRepository.Withdraw(ctx, userID, orderNumber, sum)
}

func (s *BalanceServiceImpl) GetWithdrawals(ctx context.Context, userID int) ([]models.WithdrawalModel, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	withdrawals, err := s.balanceRepository.GetWithdrawals(ctx, userID)

	var balanceRepositoryError repositories.BalanceRepositoryError
	if errors.As(err, &balanceRepositoryError) {
		switch balanceRepositoryError.Code {
		case repositories.BalanceRepositoryErrorUserNotFound:
			return models.EmptyArrayOfWithdrawalModel, NewBalanceServiceError(BalanceServiceErrorUserNotFound, "User not found")
		}
	}

	if err != nil {
		return models.EmptyArrayOfWithdrawalModel, fmt.Errorf("failed to get withdrawals: %w", err)
	}

	return withdrawals, nil
}
