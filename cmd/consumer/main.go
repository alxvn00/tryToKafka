package main

import (
	"context"
	"go.uber.org/zap"
	"myprojectKafka/internal/config"
	"myprojectKafka/internal/handler"
	"myprojectKafka/internal/infrastructure/kafka"
	"myprojectKafka/internal/infrastructure/logger"
	"os"
	"os/signal"
	"syscall"
)

const (
	topic         = "mytopic"
	consumerGroup = "my-consumer-group"
)

func main() {
	logger.Init()
	l := logger.GetLogger()
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		l.Fatal("Error loading config", zap.Error(err))
	}

	h := handler.NewMessageHandler(l, ctx)
	c, err := kafka.NewConsumer(h, cfg.Addresses, topic, consumerGroup, l)
	if err != nil {
		l.Error("error creating consumer", err)
	}

	go func() {
		c.Start()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	l.Fatal(c.Stop())
}
