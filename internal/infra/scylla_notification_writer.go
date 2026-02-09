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

func (w *ScyllaWriter) SaveNotification(ctx context.Context, n domain.NotificationDB) (domain.NotificationDB, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := w.session.Query(
		`INSERT INTO notifications (user_id, id, created_at, kind, title, body, category, deep_link, data, read) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		n.UserID, n.ID, n.CreatedAt, n.Kind, n.Title, n.Body, n.Category, n.DeepLink, n.Data, n.Read,
	).WithContext(ctx).Exec()

	if err != nil {
		return domain.NotificationDB{}, err
	}

	return n, nil
}

func (w *ScyllaWriter) UpdateNotification(ctx context.Context, userID string, id string, fields map[string]interface{}) (domain.NotificationDB, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	setClauses := ""
	values := make([]interface{}, 0, len(fields)+2)
	for col, val := range fields {
		if setClauses != "" {
			setClauses += ", "
		}
		setClauses += col + " = ?"
		values = append(values, val)
	}
	values = append(values, userID, id)

	err := w.session.Query(
		`UPDATE notifications SET `+setClauses+` WHERE user_id = ? AND id = ?`,
		values...,
	).WithContext(ctx).Exec()

	if err != nil {
		return domain.NotificationDB{}, err
	}

	return w.GetNotification(ctx, id, userID)
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
