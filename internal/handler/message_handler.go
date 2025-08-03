package handler

import (
	"encoding/json"
	"fmt"
	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"myprojectKafka/internal/entity"
)

func (h *HandlerImpl) HandleMessage(message []byte, offset confluentKafka.Offset) error {
	h.logger.Infof("Message from kafka with offset %v, %s", offset, string(message))

	var event entity.ActivityEvent
	err := json.Unmarshal(message, &event)
	if err != nil {
		h.logger.Error("Error unmarshalling message", zap.Error(err))
		return fmt.Errorf("error unmarshalling message: %w", err)
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		h.logger.Error("Error marshalling message to JSON", zap.Error(err))
		return fmt.Errorf("error marshaling massage to JSON: %w", err)
	}

	h.logger.Infof("Message from kafka with offset %v, %s", offset, string(eventJSON))

	return nil
}
