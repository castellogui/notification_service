package interfaces

import (
	"context"
	"notification_service/internal/pusher/domain"
)

type Writer interface {
	SaveNotification(ctx context.Context, n domain.NotificationDB) (domain.NotificationDB, error)
	GetNotification(ctx context.Context, id string, userID string) (domain.NotificationDB, error)
	UpdateNotification(ctx context.Context, userID string, id string, fields map[string]interface{}) (domain.NotificationDB, error)
	DeleteNotification(ctx context.Context, id string) error
}

// TODO: Remove this from pusher cause is used also in API

