package generator

import (
	"encoding/json"
	"github.com/google/uuid"
	"math/rand"
	"myprojectKafka/internal/entity"
	"time"
)

const (
	maxUserIDNumber = 100_000
)

type Generator struct {
	eventTypes []string
}

func NewGenerator() *Generator {
	return &Generator{
		eventTypes: []string{"click", "view", "like", "purchase", "login", "logout"},
	}
}

func (g *Generator) Generate() ([]byte, error) {
	event := entity.ActivityEvent{
		EventID:   uuid.New().String(),
		UserID:    rand.Intn(maxUserIDNumber),
		EventType: g.eventTypes[rand.Intn(len(g.eventTypes))],
		Timestamp: time.Now().UTC(),
	}

	return json.Marshal(&event)
}
