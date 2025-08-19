package data

type DeviceV1 struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Model      string `json:"model"`
	OrgId      string `json:"org_id"`
	Enabled    bool   `json:"enabled"`
	MacAddress string `json:"mac_address"`
	IpAddress  string `json:"ip_address"`
}

func (b DeviceV1) Clone() DeviceV1 {
	return DeviceV1{
		Id:         b.Id,
		Type:       b.Type,
		Name:       b.Name,
		Model:      b.Model,
		OrgId:      b.OrgId,
		Enabled:    b.Enabled,
		MacAddress: b.MacAddress,
		IpAddress:  b.IpAddress,
	}
}
