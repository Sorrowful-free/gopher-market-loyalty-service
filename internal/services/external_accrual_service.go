package services

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

//go:generate mockgen -source=external_accrual_service.go -destination=mock_external_accrual_service.go -package=services
type ExternalAccrualService interface {
	GetScroing(ctx context.Context, orderNumber string) (float64, error)
}

type ExternalAccrualServiceImpl struct {
	accrualSystemAddress string
	client               *resty.Client
}

func NewExternalAccrualService(accrualSystemAddress string) ExternalAccrualService {
	client := resty.New()
	return &ExternalAccrualServiceImpl{accrualSystemAddress: accrualSystemAddress, client: client}
}

func (s *ExternalAccrualServiceImpl) GetScroing(ctx context.Context, orderNumber string) (float64, error) {

	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	url, err := url.JoinPath(s.accrualSystemAddress, GetOrderPath, orderNumber)
	if err != nil {
		return 0, NewExternalAccrualServiceError(ExternalAccrualServiceErrorInternalError, "Internal server error")
	}

	resp, err := s.client.R().
		SetContext(ctx).
		Get(url)
	if err != nil {
		return 0, NewExternalAccrualServiceError(ExternalAccrualServiceErrorInternalError, "Internal server error")
	}
	if code := resp.StatusCode(); code != 200 {
		switch code {
		case 204:
			return 0, NewExternalAccrualServiceError(ExternalAccrualServiceErrorOrderNotRegistered, "Order not registered")
		case 429:
			return 0, NewExternalAccrualServiceError(ExternalAccrualServiceErrorOrderTooManyRequests, "Order too many requests")
		case 500:
			return 0, NewExternalAccrualServiceError(ExternalAccrualServiceErrorInternalError, "Internal server error")
		}
	}

	var getOrderResponse models.GetOrderResponse
	err = json.Unmarshal(resp.Body(), &getOrderResponse)
	if err != nil {
		return 0, NewExternalAccrualServiceError(ExternalAccrualServiceErrorInternalError, "Internal server error")
	}
	return 0, nil
}
