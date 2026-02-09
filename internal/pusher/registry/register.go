package registry

import (
	"encoding/json"
	"fmt"

	"notification_service/internal/pusher/domain"
)

type Decoder func(raw json.RawMessage) (domain.KindVersionMetadata, error)

type Builder func(p domain.KindVersionMetadata) (domain.ViewModel, error)

type Registry struct {
	decoders map[domain.Kind]map[int]Decoder
	builders map[domain.Kind]Builder
}

func New() *Registry {
	return &Registry{
		decoders: make(map[domain.Kind]map[int]Decoder),
		builders: make(map[domain.Kind]Builder),
	}
}

func (r *Registry) RegisterDecoder(kind domain.Kind, version int, d Decoder) {
	if _, ok := r.decoders[kind]; !ok {
		r.decoders[kind] = make(map[int]Decoder)
	}
	if _, exists := r.decoders[kind][version]; exists {
		panic(fmt.Sprintf("decoder already registered: %s v%d", kind, version))
	}
	r.decoders[kind][version] = d
}

func (r *Registry) RegisterBuilder(kind domain.Kind, b Builder) {
	if _, exists := r.builders[kind]; exists {
		panic(fmt.Sprintf("builder already registered: %s", kind))
	}
	r.builders[kind] = b
}

func (r *Registry) Decode(kind domain.Kind, version int, raw json.RawMessage) (domain.KindVersionMetadata, error) {
	byVer, ok := r.decoders[kind]
	if !ok {
		return nil, fmt.Errorf("unknown kind: %s", kind)
	}
	d, ok := byVer[version]
	if !ok {
		return nil, fmt.Errorf("unknown version %d for kind %s", version, kind)
	}
	return d(raw)
}

func (r *Registry) Build(kind domain.Kind, p domain.KindVersionMetadata) (domain.ViewModel, error) {
	b, ok := r.builders[kind]
	if !ok {
		return domain.ViewModel{}, fmt.Errorf("no builder for kind %s", kind)
	}
	return b(p)
}
