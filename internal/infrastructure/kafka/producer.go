package kafka

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"myprojectKafka/internal/config"
	"strings"
)

const (
	flushTimeout = 5000 //ms
)

var (
	errUnknownType = errors.New("unknown event type")
)

type Producer struct {
	producer *kafka.Producer
	cfg      *config.Config
	logger   *zap.SugaredLogger
}

func NewProducer(cfg *config.Config, l *zap.SugaredLogger) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(cfg.Addresses, ","),
	}
	p, err := kafka.NewProducer(conf)
	if err != nil {
		l.Error("Failed to create producer", zap.Error(err))
		return nil, fmt.Errorf("error creating new kafka producer: %w", err)
	}

	return &Producer{
		producer: p,
		cfg:      cfg,
		logger:   l,
	}, nil
}

func (p *Producer) Produce(message []byte, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
		Key:   nil,
	}

	kafkaChan := make(chan kafka.Event)

	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		p.logger.Error("Failed to produce message", zap.Error(err))
		return fmt.Errorf("error create produce message: %w", err)
	}

	e := <-kafkaChan

	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.logger.Info("producer flush-close")
	p.producer.Close()
}
