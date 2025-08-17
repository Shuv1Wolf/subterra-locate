package events

import "time"

type DevicePositioningEventV1 struct {
	DeviceId string    `json:"device_id"`
	X        float64   `json:"x"`
	Y        float64   `json:"y"`
	Z        float64   `json:"z"`
	Time     time.Time `json:"time"`
}
