package pusher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"notification_service/internal/pusher/domain"
	"notification_service/internal/pusher/registry"
	"notification_service/internal/pusher/interfaces"
)

type Handler struct {
	reg  *registry.Registry
	dbWriter interfaces.Writer
}

func NewHandler(dbWriter interfaces.Writer, reg *registry.Registry) Handler {
	return Handler{dbWriter: dbWriter, reg: reg}
}

func (h Handler) HandleMessage(ctx context.Context, raw []byte, to domain.Recipient) error {
	var env domain.Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return fmt.Errorf("unmarshal envelope: %w", err)
	}

	payload, err := h.reg.Decode(env.Kind, env.Version, env.Payload)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	vm, err := h.reg.Build(env.Kind, payload)
	if err != nil {
		return fmt.Errorf("build view model: %w", err)
	}

	err = h.dbWriter.SaveNotification(ctx, domain.NotificationDB{
		UserID: env.UserID,
		CreatedAt: *env.CreatedAt,
		ID: env.ID,
		Kind: string(env.Kind),
		Title: vm.Title,
		Body: vm.Body,
		Category: vm.Category,
		DeepLink: vm.DeepLink,
		Data: vm.Data,
		Read: false,
	})

	if err == nil {
		log.Println("notification hanlded successfully")
	}

	return err
}
