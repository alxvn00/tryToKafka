package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"myprojectKafka/internal/handler"
	"strings"
	"time"
)

const (
	noTimeout       = -1
	maxRetries      = 3
	processingDelay = 500 * time.Millisecond
)

type Consumer struct {
	cunsumer *kafka.Consumer
	handler  handler.MessageHandler
	logger   *zap.SugaredLogger
	stop     bool
}

func NewConsumer(handler handler.MessageHandler, address []string, topic, consumerGroup string, logger *zap.SugaredLogger) (*Consumer, error) {
	cfg := kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 consumerGroup,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       false,
		"auto.offset.reset":        "earliest",
	}

	c, err := kafka.NewConsumer(&cfg)
	if err != nil {
		return nil, fmt.Errorf("Error creating consumer: %w", err)
	}

	err = c.Subscribe(topic, nil)
	if err != nil {
		return nil, fmt.Errorf("Error subscribing to topic %s: %w", topic, err)
	}

	return &Consumer{
		cunsumer: c,
		handler:  handler,
		logger:   logger,
	}, nil
}

func (c *Consumer) Start() {
	for {
		if c.stop {
			break
		}
		// Чтение из кафки
		kafkaMsg, err := c.cunsumer.ReadMessage(noTimeout)
		if err != nil {
			c.logger.Error("Error reading message", zap.Error(err))
		}
		// Пропуск пустых сообщений
		if kafkaMsg == nil {
			continue
		}

		// Цикл на повторную прочитку
		retries := 0
		for retries < maxRetries {
			err = c.handler.HandleMessage(kafkaMsg.Value, kafkaMsg.TopicPartition.Offset)
			if err != nil {
				c.logger.Warnf("Error processing message, retry %d: %v\n", retries+1, err)
				retries++
				time.Sleep(processingDelay)
				continue
			}

			if _, err = c.cunsumer.StoreMessage(kafkaMsg); err != nil {
				c.logger.Error("Error storing massage", zap.Error(err))
				continue
			}

			_, err = c.cunsumer.CommitMessage(kafkaMsg)
			if err != nil {
				c.logger.Error("Error committing message", zap.Error(err))
			}

			break
		}
		// Если сообщение не удалось прочитать логирование ошибки и сообщения
		if retries == maxRetries {
			c.logger.Warn("Failed to process message after %d retries: %s\n", retries, string(kafkaMsg.Value))
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	if _, err := c.cunsumer.Commit(); err != nil {
		c.logger.Error("Error committing message", zap.Error(err))
		return err
	}
	return c.cunsumer.Close()
}
