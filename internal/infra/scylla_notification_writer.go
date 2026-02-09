package infra

import (
	"context"
	"time"

	"github.com/gocql/gocql"

	"notification_service/internal/pusher/domain"
	"notification_service/internal/pusher/interfaces"
)

type ScyllaWriter struct {
	session *gocql.Session
}

func NewScyllaWriter(session *gocql.Session) interfaces.Writer {
	return &ScyllaWriter{session: session}
}

func (w *ScyllaWriter) SaveNotification(ctx context.Context, n domain.NotificationDB) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return w.session.Query(
		`INSERT INTO notifications (user_id, created_at, id, kind, title, body, category, deep_link, data, read) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		n.UserID, n.CreatedAt, n.ID, n.Kind, n.Title, n.Body, n.Category, n.DeepLink, n.Data, n.Read,
	).WithContext(ctx).Exec()
}
