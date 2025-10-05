package app

import (
	"database/sql"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/config"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/handlers"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
)

type App struct {
	logger logger.Logger
	config config.Config

	db *sql.DB

	userRepository repositories.UserRepository
	userService    services.UserService

	handlers handlers.Handlers
}

func NewApp() *App {
	return &App{}
}

func (a *App) BuildConfig() error {
	a.config = config.NewLocalConfig()
	err := a.config.Parse()
	if err != nil {
		a.logger.Error("Failed to parse config", "error", err)
		return err
	}
	return nil
}

func (a *App) BuildLogger() error {
	a.logger = logger.NewZapLogger()
	return nil
}

func (a *App) BuildDatabase() error {
	db, err := sql.Open("postgres", a.config.DatabaseURI())
	if err != nil {
		a.logger.Error("Failed to open database", "error", err)
		return err
	}
	a.db = db
	return nil
}

func (a *App) BuildRepositories() error {
	a.userRepository = repositories.NewPGUserRepository(a.db)
	return nil
}

func (a *App) BuildServices() error {
	a.userService = services.NewUserService(a.userRepository)
	return nil
}

func (a *App) BuildHandlers() error {
	a.handlers = handlers.NewEchoHandlers(a.logger)
	a.handlers.BuildGroups()
	a.handlers.BuildAuthMiddleware(a.config.JwtSecret())
	a.handlers.BuildRoutes()
	return nil
}

func (a *App) Run() error {
	return a.handlers.Run()
}
