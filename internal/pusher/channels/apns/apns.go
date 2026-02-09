package apns

import (
	"context"
	"log"

	"notification_service/internal/pusher/domain"
)

type Adapter struct{}

func NewAdapter() Adapter {
	return Adapter{}
}

func (a Adapter) Send(ctx context.Context, to domain.Recipient, vm domain.ViewModel) error {
	log.Printf("APNs: sending to %s: %s - %s\n", to.DeviceToken, vm.Title, vm.Body)
	return nil
}
