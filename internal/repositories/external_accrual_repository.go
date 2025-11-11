package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/go-resty/resty/v2"
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
	client               *resty.Client
}

func NewExternalAccrualRepository(accrualSystemAddress string) ExternalAccrualRepository {
	client := resty.New()
	return &ExternalAccrualRepositoryImpl{accrualSystemAddress: accrualSystemAddress, client: client}
}

func (r *ExternalAccrualRepositoryImpl) GetScoring(ctx context.Context, orderNumber string) (models.ScoringModel, error) {

	if ctx.Err() != nil {
		return models.EmptyScoringModel, ctx.Err()
	}

	url, err := url.JoinPath(r.accrualSystemAddress, GetOrderPath, orderNumber)
	if err != nil {
		return models.EmptyScoringModel, fmt.Errorf("failed to join path: %w", err)
	}

	resp, err := r.client.R().
		SetContext(ctx).
		Get(url)
	if err != nil {
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

	var scoringModel models.ScoringModel
	err = json.Unmarshal(resp.Body(), &scoringModel)
	if err != nil {
		return models.EmptyScoringModel, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorInternalError, "Internal server error")
	}
	return scoringModel, nil
}
