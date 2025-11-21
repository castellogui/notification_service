package notifications

import "time"

type NotificationEnvelope struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}
