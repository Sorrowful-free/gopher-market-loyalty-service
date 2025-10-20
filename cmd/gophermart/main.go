package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/app"
)

func main() {

	app := app.NewApp()
	if err := app.BuildConfig(); err != nil {
		log.Fatalf("Failed to build config: %v", err)
	}
	if err := app.BuildLogger(); err != nil {
		log.Fatalf("Failed to build logger: %v", err)
	}
	if err := app.BuildDatabase(); err != nil {
		log.Fatalf("Failed to build database: %v", err)
	}
	if err := app.BuildRepositories(); err != nil {
		log.Fatalf("Failed to build repositories: %v", err)
	}
	if err := app.BuildServices(); err != nil {
		log.Fatalf("Failed to build services: %v", err)
	}
	if err := app.BuildHandlers(); err != nil {
		log.Fatalf("Failed to build handlers: %v", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
