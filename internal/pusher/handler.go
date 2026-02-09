package pusher

import (
	"context"
	"encoding/json"
	"fmt"

	"notification_service/internal/pusher/channels/apns"
	"notification_service/internal/pusher/domain"
	"notification_service/internal/pusher/registry"
)

type Handler struct {
	apns apns.Adapter
	reg  *registry.Registry
}

func NewHandler(a apns.Adapter, reg *registry.Registry) Handler {
	return Handler{apns: a, reg: reg}
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

	return h.apns.Send(ctx, to, vm)
}
