package repositories

import (
	context "context"
	"errors"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type PGXUserRepository struct {
	pgxPool *pgxpool.Pool
}

func NewPGXUserRepository(pgxPool *pgxpool.Pool) UserRepository {
	return &PGXUserRepository{pgxPool: pgxPool}
}

func (r *PGXUserRepository) Create(ctx context.Context, login string, password string) (models.UserModel, error) {
	const query = `
		INSERT INTO users (login, pass_hash)
		VALUES ($1, $2)
		RETURNING id, login, pass_hash
	`

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.EMPTY_USER_MODEL, NewUserRepositoryError(UserRepositoryErrorInternalError, "Failed to generate password hash")
	}
	var user models.UserModel
	err = r.pgxPool.QueryRow(ctx, query, login, string(hash)).Scan(&user.ID, &user.Login, &user.Password)

	var pgxErr *pgconn.PgError
	if errors.As(err, &pgxErr) {
		if pgxErr.Code == "23505" {
			return models.EMPTY_USER_MODEL, NewUserRepositoryError(UserRepositoryErrorUserAlreadyExists, "User already exists")
		}
	}

	if err != nil {
		return models.EMPTY_USER_MODEL, NewUserRepositoryError(UserRepositoryErrorInternalError, "Failed to create user")
	}
	return user, nil
}

func (r *PGXUserRepository) GetByLoginAndPassword(ctx context.Context, login string, password string) (models.UserModel, error) {

	const query = `
		SELECT id, login, pass_hash
		FROM users
		WHERE login = $1
	`
	row := r.pgxPool.QueryRow(ctx, query, login)
	var user models.UserModel
	var passHash string
	err := row.Scan(&user.ID, &user.Login, &passHash)
	if err != nil && err != pgx.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password)) != nil {
		return models.EMPTY_USER_MODEL, NewUserRepositoryError(UserRepositoryErrorInvalidCredentials, "Invalid credentials")
	}
	if err == pgx.ErrNoRows {
		return models.EMPTY_USER_MODEL, NewUserRepositoryError(UserRepositoryErrorUserNotFound, "User not found")
	}
	return user, nil
}

func (r *PGXUserRepository) GetByID(ctx context.Context, id int) (models.UserModel, error) {

	const query = `
		SELECT id, login, password
		FROM users
		WHERE id = $1
	`
	row := r.pgxPool.QueryRow(ctx, query, id)
	var user models.UserModel
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil && err != pgx.ErrNoRows {
		return models.EMPTY_USER_MODEL, NewUserRepositoryError(UserRepositoryErrorInternalError, "Failed to get user by id")
	}
	if err == pgx.ErrNoRows {
		return models.EMPTY_USER_MODEL, NewUserRepositoryError(UserRepositoryErrorUserNotFound, "User not found by id")
	}
	return user, nil
}

func (r *PGXUserRepository) ComparePassword(ctx context.Context, password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
