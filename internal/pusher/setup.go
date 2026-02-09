package pusher

import (
	"notification_service/internal/pusher/kinds"
	"notification_service/internal/pusher/registry"
)

func SetupRegistry() *registry.Registry {
	reg := registry.New()
	kinds.StatusRegister(reg)
	return reg
}
