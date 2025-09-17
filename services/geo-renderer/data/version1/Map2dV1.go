package data1

import "time"

type Map2dV1 struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	SVG       string    `json:"svg_content"`
	ScaleX    float64   `json:"scale_x"`
	ScaleY    float64   `json:"scale_y"`
	CreatedAt time.Time `json:"created_at"`
	OrgId     string    `json:"org_id"`
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
