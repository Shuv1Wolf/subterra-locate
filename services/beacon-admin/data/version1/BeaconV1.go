package data

type BeaconV1 struct {
	Id      string  `json:"id"`
	Type    string  `json:"type"`
	Udi     string  `json:"udi"`
	Label   string  `json:"label"`
	X       float32 `json:"x"`
	Y       float32 `json:"y"`
	Z       float32 `json:"z"`
	OrgId   string  `json:"org_id"`
	Enabled bool    `json:"enabled"`
	MapId   string  `json:"mao_id"`
	// Radius float32    `json:"radius" bson:"radius"`
}

func (b BeaconV1) Clone() BeaconV1 {
	return BeaconV1{
		Id:      b.Id,
		Type:    b.Type,
		Udi:     b.Udi,
		Label:   b.Label,
		X:       b.X,
		Y:       b.Y,
		Z:       b.Z,
		OrgId:   b.OrgId,
		Enabled: b.Enabled,
		MapId:   b.MapId,
		// Radius: b.Radius,
	}
}
