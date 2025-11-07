package main

import (
	"context"
	"log"

	"app/internal/app"
)

func main() {
	a := app.New()

	ctx := context.Background()
	err := a.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
