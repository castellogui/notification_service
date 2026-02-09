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
		`INSERT INTO notifications (user_id,  id, created_at, kind, title, body, category, deep_link, data, read) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		n.UserID,n.ID, n.CreatedAt, n.Kind, n.Title, n.Body, n.Category, n.DeepLink, n.Data, n.Read,
	).WithContext(ctx).Exec()
}

func (w *ScyllaWriter) DeleteNotification(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return w.session.Query(
		`DELETE FROM notifications WHERE id = ?`,
		id,
	).WithContext(ctx).Exec()
}

func (w *ScyllaWriter) GetNotification(ctx context.Context, id string, userID string) (domain.NotificationDB, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var n domain.NotificationDB
	err := w.session.Query(
		`SELECT user_id, id, created_at, kind, title, body, category, deep_link, data, read FROM notifications WHERE user_id = ? AND id = ?`,
		userID, id,
	).WithContext(ctx).Scan(&n.UserID, &n.ID, &n.CreatedAt, &n.Kind, &n.Title, &n.Body, &n.Category, &n.DeepLink, &n.Data, &n.Read)

	if err != nil {
		return domain.NotificationDB{}, err
	}

	return n, nil
}
