package repositories

import (
	"database/sql"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
)

type PGUserRepository struct {
	db *sql.DB
}

func NewPGUserRepository(db *sql.DB) UserRepository {
	return &PGUserRepository{db: db}
}

func (r *PGUserRepository) Create(login string, password string) (models.UserModel, error) {
	return models.UserModel{}, nil
}

func (r *PGUserRepository) GetByLoginAndPassword(login string, password string) (models.UserModel, error) {
	return models.UserModel{}, nil
}

func (r *PGUserRepository) GetByID(id int64) (models.UserModel, error) {
	return models.UserModel{}, nil
}
