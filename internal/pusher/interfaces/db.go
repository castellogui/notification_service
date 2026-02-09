package interfaces

import (
	"context"
	"notification_service/internal/pusher/domain"
)

type Writer interface {
	SaveNotification(ctx context.Context, n domain.NotificationDB) error
}

