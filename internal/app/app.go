package app

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"myprojectKafka/internal/config"
	"myprojectKafka/internal/generator"
	"myprojectKafka/internal/infrastructure/kafka"
	"myprojectKafka/internal/infrastructure/logger"
	"myprojectKafka/internal/service"
	"os"
)

func Start(ctx context.Context) error {
	logger.Init()
	l := logger.GetLogger()

	defer func() {
		if err := l.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			fmt.Fprintf(os.Stderr, "Error closing logger: %v\n", err)
		}
	}()

	if err := godotenv.Load(); err != nil {
		l.Warn("Error loading .env file", zap.Error(err))
		return fmt.Errorf("error loading .env file: %w", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		l.Fatal("Config error", zap.Error(err))
		return fmt.Errorf("app.Start config error: %w", err)
	}

	prod, err := kafka.NewProducer(cfg, l)
	if err != nil {
		l.Errorf("Failed to create kafka producer", zap.Error(err))
	}
	defer prod.Close()

	gen := generator.NewGenerator()

	src := service.NewService(gen, prod, l)
	if err = src.PushRandomData(); err != nil {
		l.Errorf("Failed to push data to kafka producer", zap.Error(err))
		return fmt.Errorf("app.Start random data service error: %w", err)
	}

	return nil
}
