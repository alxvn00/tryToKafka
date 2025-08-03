package main

import (
	"context"
	"log"
	"myprojectKafka/internal/app"
)

func main() {
	ctx := context.Background()

	if err := app.Start(ctx); err != nil {
		log.Fatalf("main not started %w", err)
	}
}
