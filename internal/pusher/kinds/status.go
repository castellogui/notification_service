package kinds

import (
	"encoding/json"
	"fmt"

	"notification_service/internal/pusher/domain"
	"notification_service/internal/pusher/registry"
)

type PayloadV1 struct {
	Entity    string `json:"entity"`
	EntityID  string `json:"entity_id"`
	NewStatus string `json:"new_status"`
}

func (s PayloadV1) Kind() domain.Kind   { return domain.Kind("status") }
func (s PayloadV1) Version() int        { return 1 }
func (s PayloadV1) Validate() error     { return nil }

func decodeV1(raw json.RawMessage) (domain.KindVersionMetadata, error) {
	var s PayloadV1
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, fmt.Errorf("error unmarshalling status v1: %w", err)
	}
	return s, s.Validate()
}

func build(p domain.KindVersionMetadata) (domain.ViewModel, error) {
	v, ok := p.(PayloadV1)
	if !ok {
		return domain.ViewModel{}, fmt.Errorf("status builder expects PayloadV1")
	}
	return domain.ViewModel{
		Title:    "Atualização",
		Body:     v.Entity + " " + v.EntityID + " -> " + v.NewStatus,
		Category: "status",
		DeepLink: "app://orders/" + v.EntityID,
		Data: map[string]string{"entity_id": v.EntityID, "new_status": v.NewStatus},
	}, nil
}

func StatusRegister(reg *registry.Registry) {
	reg.RegisterDecoder(domain.Kind("status"), 1, decodeV1)
	reg.RegisterBuilder(domain.Kind("status"), build)
}
