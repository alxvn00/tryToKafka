package handler

import (
	"context"
	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type MessageHandler interface {
	HandleMessage(message []byte, offset confluentKafka.Offset) error
}

type HandlerImpl struct {
	logger *zap.SugaredLogger
	ctx    context.Context
}

func NewMessageHandler(l *zap.SugaredLogger, ctx context.Context) MessageHandler {
	return &HandlerImpl{
		logger: l,
		ctx:    ctx,
	}
}
