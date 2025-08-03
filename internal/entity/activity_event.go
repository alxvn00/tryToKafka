package entity

import "time"

type ActivityEvent struct {
	EventID   string    `json:"event_id"`
	UserID    int       `json:"user_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}
