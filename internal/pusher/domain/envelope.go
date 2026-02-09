package domain

import (
	"encoding/json"
	"time"
)

type Kind string
type Channel string

type Envelope struct {
	ID				string			`json:"id"`
	UserID			string  		`json:"userid"`
	Kind			Kind			`json:"kind"`
	Version 		int				`json:"version"`
	DeliverAt		*time.Time		`json:"deliverAt,omitempty"`
	ChannelHints 	[]Channel		`json:"channelHints,omitempty"`
	Payload			json.RawMessage	`json:"payload"`
}

type KindVersionMetadata interface {
	Kind() Kind
	Version() int
	Validate() error
}

type Recipient struct {
	DeviceToken string
}