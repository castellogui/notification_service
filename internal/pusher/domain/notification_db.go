package domain

import "time"

// NotificationDB é o modelo de persistência de notificação (tabela notifications).
type NotificationDB struct {
	UserID    string
	CreatedAt time.Time
	ID        string
	Kind      string
	Title     string
	Body      string
	Category  string
	DeepLink  string
	Data      map[string]string
	Read      bool
}
