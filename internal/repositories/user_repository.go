package repositories

import "github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"

type UserRepository interface {
	Create(login string, password string) (models.UserModel, error)
	GetByLoginAndPassword(login string, password string) (models.UserModel, error)

	GetBalance(userID string) (models.BalanceModel, error)
}
