package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeviceStatus string

const (
	DeviceStatusUnknown DeviceStatus = "unknown"
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusIdle    DeviceStatus = "idle"
	DeviceStatusOffline DeviceStatus = "offline"
)

type Device struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name              string             `bson:"name"`
	IPv4              string             `bson:"ipv4"`
	MAC               *string            `bson:"mac,omitempty"`
	Vendor            *string            `bson:"vendor,omitempty"`
	Status            DeviceStatus       `bson:"status"`
	NetworkID         primitive.ObjectID `bson:"network_id,omitempty"`
	Ports             []Port             `bson:"ports,omitempty"`
	Hostname          *string            `bson:"hostname,omitempty"`
	CreatedAt         time.Time          `bson:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at"`
	LastSeenOnlineAt  *time.Time         `bson:"last_seen_online_at,omitempty"`
	PortScanStartedAt *time.Time         `bson:"port_scan_started_at,omitempty"`
	PortScanEndedAt   *time.Time         `bson:"port_scan_ended_at,omitempty"`
}
