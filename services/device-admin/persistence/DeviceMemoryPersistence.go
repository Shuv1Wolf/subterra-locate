package persistence

import (
	"context"
	"fmt"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
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

func (c *DeviceMemoryPersistence) GetPageByFilter(ctx context.Context, reqctx cdata.RequestContextV1,
	filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.DeviceV1], error) {

	if reqctx.OrgId != "" {
		filter.Put("org_id", reqctx.OrgId)
	}

	return c.IdentifiableMemoryPersistence.
		GetPageByFilter(ctx, c.composeFilter(filter), paging, nil, nil)
}

func (c *DeviceMemoryPersistence) GetOneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.DeviceV1, error) {
	item, err := c.IdentifiableMemoryPersistence.GetOneById(ctx, id)
	if err != nil {
		return item, err
	}

	if item.Id != "" && item.OrgId != reqctx.OrgId {
		return data.DeviceV1{}, nil
	}

	return item, nil
}

func (c *DeviceMemoryPersistence) Create(ctx context.Context, reqctx cdata.RequestContextV1, item data.DeviceV1) (data.DeviceV1, error) {
	if item.OrgId == "" {
		item.OrgId = reqctx.OrgId
	}
	return c.IdentifiableMemoryPersistence.Create(ctx, item)
}

func (c *DeviceMemoryPersistence) Update(ctx context.Context, reqctx cdata.RequestContextV1, item data.DeviceV1) (data.DeviceV1, error) {
	data, err := c.GetOneById(ctx, reqctx, item.Id)
	if err != nil {
		return data, err
	}

	if data.Id == "" {
		return data, fmt.Errorf("device not found: %s", item.Id)
	}

	return c.IdentifiableMemoryPersistence.Update(ctx, item)
}

func (c *DeviceMemoryPersistence) DeleteById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.DeviceV1, error) {
	data, err := c.GetOneById(ctx, reqctx, id)
	if err != nil {
		return data, err
	}

	if data.Id == "" {
		return data, fmt.Errorf("device not found: %s", id)
	}

	return c.IdentifiableMemoryPersistence.DeleteById(ctx, id)
}
