package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const (
	GetOrderPath = "/api/orders/"
)

//go:generate mockgen -source=external_accrual_repository.go -destination=mock_external_accrual_repository.go -package=repositories
type ExternalAccrualRepository interface {
	GetScoring(ctx context.Context, orderNumber string) (models.ScoringModel, error)
}

type ExternalAccrualRepositoryImpl struct {
	accrualSystemAddress string
	logger               logger.Logger
	client               *resty.Client
}

func NewExternalAccrualRepository(accrualSystemAddress string, logger logger.Logger) ExternalAccrualRepository {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	client := resty.New().SetLogger(zapLogger.Sugar())
	return &ExternalAccrualRepositoryImpl{accrualSystemAddress: accrualSystemAddress, logger: logger, client: client}
}

func (r *ExternalAccrualRepositoryImpl) GetScoring(ctx context.Context, orderNumber string) (models.ScoringModel, error) {

	if ctx.Err() != nil {
		r.logger.Error("context error", "error", ctx.Err())
		return models.EmptyScoringModel, ctx.Err()
	}

	url, err := url.JoinPath(r.accrualSystemAddress, GetOrderPath, orderNumber)
	if err != nil {
		r.logger.Error("failed to join path", "error", err)
		return models.EmptyScoringModel, fmt.Errorf("failed to join path: %w", err)
	}

	resp, err := r.client.R().
		SetContext(ctx).
		Get(url)
	if err != nil {
		r.logger.Error("failed to get scoring", "error", err)
		return models.EmptyScoringModel, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorInternalError, "Internal server error")
	}
	if code := resp.StatusCode(); code != 200 {
		switch code {
		case 204:
			return models.EmptyScoringModel, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorOrderNotRegistered, "Order not registered")
		case 429:
			return models.EmptyScoringModel, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorOrderTooManyRequests, "Order too many requests")
		case 500:
			return models.EmptyScoringModel, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorInternalError, "Internal server error")
		}
	}

	if resp.IsError() {
		r.logger.Error("failed to get scoring", "error", resp.Error())
		return models.EmptyScoringModel, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorInternalError, "Internal server error")
	}

	var scoringModel models.ScoringModel
	err = json.Unmarshal(resp.Body(), &scoringModel)
	if err != nil {
		r.logger.Error("failed to unmarshal scoring", "error", err)
		return models.EmptyScoringModel, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorInternalError, "Internal server error")
	}
	r.logger.Info("scoring", "scoring", scoringModel)
	return scoringModel, nil
}
