package main

import (
	"context"
	"log"

	"app/internal/app"
	"app/pkg/env"
)

func main() {
	ctx := context.Background()
	environment := env.New("")

	a, err := app.New(&app.AppConfig{
		HttpAddr: environment.Get("APP_SERVER_ADDRESS", ":8080"),

		PGUser:     environment.Get("POSTGRES_USER", "postgres_user"),
		PGPassword: environment.Get("POSTGRES_PASSWORD", "postgres_ps"),
		PGDb:       environment.Get("POSTGRES_DB", "postgres_db"),
		PGPort:     environment.Get("APP_POSTGRES_PORT", "15432"),
		PGHost:     environment.Get("APP_POSTGRES_HOST", "localhost"),
	})
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
