package domain

import "time"

type NotificationDB struct {
	UserID    string            `json:"user_id"`
	CreatedAt time.Time         `json:"created_at"`
	ID        string            `json:"id"`
	Kind      string            `json:"kind"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Category  string            `json:"category"`
	DeepLink  string            `json:"deep_link"`
	Data      map[string]string `json:"data"`
	Read      bool              `json:"read"`
}
