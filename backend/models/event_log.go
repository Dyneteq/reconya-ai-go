package models

import (
	"time"
)

type EventLog struct {
	Type        EEventLogType `bson:"type"`
	Description string        `bson:"description"`
	DeviceID    *string       `bson:"device_id,omitempty"`
	CreatedAt   *time.Time    `bson:"created_at,omitempty"`
	UpdatedAt   *time.Time    `bson:"updated_at,omitempty"`
}
