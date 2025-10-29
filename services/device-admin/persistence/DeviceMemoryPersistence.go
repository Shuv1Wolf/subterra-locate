package persistence

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type DeviceMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence[data.DeviceV1, string]
}

func NewDeviceMemoryPersistence() *DeviceMemoryPersistence {
	c := &DeviceMemoryPersistence{
		IdentifiableMemoryPersistence: *cpersist.NewIdentifiableMemoryPersistence[data.DeviceV1, string](),
	}
	c.IdentifiableMemoryPersistence.MaxPageSize = 1000
	return c
}

func (c *DeviceMemoryPersistence) composeFilter(filter cquery.FilterParams) func(beacon data.DeviceV1) bool {

	id := filter.GetAsString("id")
	OrgId := filter.GetAsString("org_id")
	name := filter.GetAsString("name")
	model := filter.GetAsString("model")
	macAddress := filter.GetAsString("mac_address")

	return func(beacon data.DeviceV1) bool {
		if id != "" && beacon.Id != id {
			return false
		}
		if OrgId != "" && beacon.OrgId != OrgId {
			return false
		}
		if name != "" && beacon.Name != name {
			return false
		}
		if model != "" && beacon.Model != model {
			return false
		}
		if macAddress != "" && beacon.MacAddress != macAddress {
			return false
		}
		return true
	}
}

func (c *DeviceMemoryPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.DeviceV1], error) {

	return c.IdentifiableMemoryPersistence.
		GetPageByFilter(ctx, c.composeFilter(filter), paging, nil, nil)
}

func ContainsStr(arr []string, substr string) bool {
	for _, _str := range arr {
		if _str == substr {
			return true
		}
	}
	return false
}
