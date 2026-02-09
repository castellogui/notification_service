package setup

import (
	"notification_service/internal/pusher/registry"
	"notification_service/internal/pusher/kinds"
)

func SetupRegistry() *registry.Registry {
	reg := registry.New()
	kinds.StatusRegister(reg)
	return reg
}