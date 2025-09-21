package data1

import "time"

type Map2dV1 struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	SVG       string    `json:"svg_content"` // original svg
	ScaleX    float64   `json:"scale_x"`     // scale in meters
	ScaleY    float64   `json:"scale_y"`     // scale in meters
	CreatedAt time.Time `json:"created_at"`
	OrgId     string    `json:"org_id"`
	Width     float64   `json:"width"`  // width in meters
	Height    float64   `json:"height"` // height in meters
	Level     int       `json:"level"`  // level of the map
}

func (m *Map2dV1) Clone() Map2dV1 {
	return Map2dV1{
		Id:        m.Id,
		Name:      m.Name,
		SVG:       m.SVG,
		ScaleX:    m.ScaleX,
		ScaleY:    m.ScaleY,
		CreatedAt: m.CreatedAt,
		OrgId:     m.OrgId,
	}
}
