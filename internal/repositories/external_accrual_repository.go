package repositories

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/go-resty/resty/v2"
)

const (
	GetOrderPath = "/api/orders/"
)

//go:generate mockgen -source=external_accrual_repository.go -destination=mock_external_accrual_repository.go -package=repositories
type ExternalAccrualRepository interface {
	GetScroing(ctx context.Context, orderNumber string) (float64, error)
}

type ExternalAccrualRepositoryImpl struct {
	accrualSystemAddress string
	client               *resty.Client
}

func NewExternalAccrualRepository(accrualSystemAddress string) ExternalAccrualRepository {
	client := resty.New()
	return &ExternalAccrualRepositoryImpl{accrualSystemAddress: accrualSystemAddress, client: client}
}

func (r *ExternalAccrualRepositoryImpl) GetScroing(ctx context.Context, orderNumber string) (float64, error) {

	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	url, err := url.JoinPath(r.accrualSystemAddress, GetOrderPath, orderNumber)
	if err != nil {
		return 0, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorInternalError, "Internal server error")
	}

	resp, err := r.client.R().
		SetContext(ctx).
		Get(url)
	if err != nil {
		return 0, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorInternalError, "Internal server error")
	}
	if code := resp.StatusCode(); code != 200 {
		switch code {
		case 204:
			return 0, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorOrderNotRegistered, "Order not registered")
		case 429:
			return 0, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorOrderTooManyRequests, "Order too many requests")
		case 500:
			return 0, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorInternalError, "Internal server error")
		}
	}

	var getOrderResponse models.GetOrderResponse
	err = json.Unmarshal(resp.Body(), &getOrderResponse)
	if err != nil {
		return 0, NewExternalAccrualRepositoryError(ExternalAccrualRepositoryErrorInternalError, "Internal server error")
	}
	return 0, nil
}
