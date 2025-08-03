package service

import (
	"go.uber.org/zap"
	"myprojectKafka/internal/generator"
	"myprojectKafka/internal/infrastructure/kafka"
	"time"
)

const (
	topic     = "mytopic"
	limitData = 100
)

var (
	timeDelayGenerate = 500 * time.Millisecond
)

type Service struct {
	gen      *generator.Generator
	producer *kafka.Producer
	l        *zap.SugaredLogger
}

func NewService(gen *generator.Generator, producer *kafka.Producer, l *zap.SugaredLogger) *Service {
	return &Service{
		gen:      gen,
		producer: producer,
		l:        l,
	}
}

func (s *Service) PushRandomData() error {
	for i := 0; i < limitData; i++ {
		msg, err := s.gen.Generate()
		if err != nil {
			s.l.Errorf("Generation error", zap.Error(err))
			continue
		}

		err = s.producer.Produce(msg, topic)
		if err != nil {
			s.l.Errorf("Produce error", zap.Error(err))
			continue
		}

		s.l.Infof("Produce success, message: %s", msg)

		time.Sleep(timeDelayGenerate)
	}
	return nil
}
