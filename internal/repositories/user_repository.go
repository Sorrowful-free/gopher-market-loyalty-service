package repositories

import "github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"

type UserRepository interface {
	Create(login string, password string) (models.UserModel, error)
	GetByLoginAndPassword(login string, password string) (models.UserModel, error)
	GetByID(id int64) (models.UserModel, error)

	CheckIfUserExists(login string) (bool, error)
	CheckIfTokenIsValid(token string) (bool, error)
}
