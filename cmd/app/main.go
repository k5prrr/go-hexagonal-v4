package main

import (
	"context"
	"log"

	"app/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
