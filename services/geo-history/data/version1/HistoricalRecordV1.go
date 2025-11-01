package data1

import "time"

type HistoricalRecordV1 struct {
	Id        string    `json:"id"`
	EntityId  string    `json:"entity_id"`
	MapId     string    `json:"map_id"`
	OrgId     string    `json:"org_id"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	Z         float64   `json:"z"`
	Timestamp time.Time `json:"timestamp"`
}
