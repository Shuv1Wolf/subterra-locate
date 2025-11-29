package data1

import "time"

type ZoneV1 struct {
	Id        string    `json:"id"`
	MapId     string    `json:"map_id"`
	OrgId     string    `json:"org_id"`
	Name      string    `json:"name"`
	PositionX float64   `json:"position_x"`
	PositionY float64   `json:"position_y"`
	Width     float64   `json:"width"`
	Height    float64   `json:"height"`
	Type      string    `json:"type"`
	Color     string    `json:"color"`
	MaxDevice int       `json:"max_device"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *ZoneV1) Clone() ZoneV1 {
	return ZoneV1{
		Id:        m.Id,
		Name:      m.Name,
		MapId:     m.MapId,
		OrgId:     m.OrgId,
		PositionX: m.PositionX,
		PositionY: m.PositionY,
		Width:     m.Width,
		Height:    m.Height,
		Type:      m.Type,
		Color:     m.Color,
		MaxDevice: m.MaxDevice,
		CreatedAt: m.CreatedAt,
	}
}
