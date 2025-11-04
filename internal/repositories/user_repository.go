package repositories

import (
	"context"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
)

//go:generate mockgen -source=user_repository.go -destination=mock_user_repository.go -package=repositories
type UserRepository interface {
	Create(ctx context.Context, login string, password string) (models.UserModel, error)
	GetByLoginAndPassword(ctx context.Context, login string, password string) (models.UserModel, error)
	GetByID(ctx context.Context, id int) (models.UserModel, error)
}
