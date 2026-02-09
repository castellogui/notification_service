package interfaces

import (
	"context"
	"notification_service/internal/pusher/domain"
)

type Writer interface {
	SaveNotification(ctx context.Context, n domain.NotificationDB) error
	GetNotification(ctx context.Context, id string, userID string) (domain.NotificationDB, error)
	DeleteNotification(ctx context.Context, id string) error
}

