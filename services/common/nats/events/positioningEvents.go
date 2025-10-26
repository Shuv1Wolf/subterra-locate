package events

import "time"

type DevicePositioningEventV1 struct {
	DeviceId  string    `json:"device_id"`
	MapId     string    `json:"map_id"`
	OrgId     string    `json:"org_id"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	Z         float64   `json:"z"`
	Timestamp time.Time `json:"timestamp"`
}
