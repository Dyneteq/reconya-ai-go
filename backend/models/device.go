package models

import (
	"time"
)

type DeviceStatus string

const (
	DeviceStatusUnknown DeviceStatus = "unknown"
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusIdle    DeviceStatus = "idle"
	DeviceStatusOffline DeviceStatus = "offline"
)

type Device struct {
	ID                string       `bson:"_id,omitempty" json:"id"`
	Name              string       `bson:"name" json:"name"`
	IPv4              string       `bson:"ipv4" json:"ipv4"`
	MAC               *string      `bson:"mac,omitempty" json:"mac,omitempty"`
	Vendor            *string      `bson:"vendor,omitempty" json:"vendor,omitempty"`
	Status            DeviceStatus `bson:"status" json:"status"`
	NetworkID         string       `bson:"network_id,omitempty" json:"network_id,omitempty"`
	Ports             []Port       `bson:"ports,omitempty" json:"ports,omitempty"`
	Hostname          *string      `bson:"hostname,omitempty" json:"hostname,omitempty"`
	CreatedAt         time.Time    `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time    `bson:"updated_at" json:"updated_at"`
	LastSeenOnlineAt  *time.Time   `bson:"last_seen_online_at,omitempty" json:"last_seen_online_at,omitempty"`
	PortScanStartedAt *time.Time   `bson:"port_scan_started_at,omitempty" json:"port_scan_started_at,omitempty"`
	PortScanEndedAt   *time.Time   `bson:"port_scan_ended_at,omitempty" json:"port_scan_ended_at,omitempty"`
}
