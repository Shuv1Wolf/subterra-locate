package persistence

import (
	"context"
	"strings"

	data "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type BeaconsMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence[data.BeaconV1, string]
}

func NewBeaconsMemoryPersistence() *BeaconsMemoryPersistence {
	c := &BeaconsMemoryPersistence{
		IdentifiableMemoryPersistence: *cpersist.NewIdentifiableMemoryPersistence[data.BeaconV1, string](),
	}
	c.MaxPageSize = 1000
	return c
}

func (c *BeaconsMemoryPersistence) composeFilter(filter cquery.FilterParams) func(item data.BeaconV1) bool {
	id, idOk := filter.GetAsNullableString("id")
	orgId, orgIdOk := filter.GetAsNullableString("org_id")
	label, labelOk := filter.GetAsNullableString("label")
	udi, udiOk := filter.GetAsNullableString("udi")
	udis, udisOk := filter.GetAsNullableString("udis")
	var udiValues []string
	if udisOk && udis != "" {
		udiValues = strings.Split(udis, ",")
	}

	return func(item data.BeaconV1) bool {
		if idOk && item.Id != id {
			return false
		}
		if orgIdOk && item.OrgId != orgId {
			return false
		}
		if labelOk && item.Label != label {
			return false
		}
		if udiOk && item.Udi != udi {
			return false
		}
		if len(udiValues) > 0 {
			found := false
			for _, v := range udiValues {
				if v == item.Udi {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
}

func (c *BeaconsMemoryPersistence) GetPageByFilter(ctx context.Context, reqctx cdata.RequestContextV1,
	filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.BeaconV1], error) {

	if reqctx.OrgId != "" {
		filter.Put("org_id", reqctx.OrgId)
	}

	return c.IdentifiableMemoryPersistence.GetPageByFilter(ctx, c.composeFilter(filter), paging, nil, nil)
}

func (c *BeaconsMemoryPersistence) GetOneByUdi(ctx context.Context, reqctx cdata.RequestContextV1, udi string) (data.BeaconV1, error) {
	var item data.BeaconV1
	found := false

	for _, beacon := range c.Items {
		if beacon.Udi == udi {
			if reqctx.OrgId != "" && beacon.OrgId != reqctx.OrgId {
				continue
			}
			item = beacon
			found = true
			break
		}
	}

	if !found {
		c.Logger.Trace(ctx, "Cannot find beacon by %s", udi)
		return data.BeaconV1{}, nil
	}

	c.Logger.Trace(ctx, "Found beacon by %s", udi)
	return item, nil
}

func (c *BeaconsMemoryPersistence) Create(ctx context.Context, reqctx cdata.RequestContextV1, item data.BeaconV1) (data.BeaconV1, error) {
	if reqctx.OrgId != "" {
		item.OrgId = reqctx.OrgId
	}
	return c.IdentifiableMemoryPersistence.Create(ctx, item)
}

func (c *BeaconsMemoryPersistence) GetOneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.BeaconV1, error) {
	item, err := c.IdentifiableMemoryPersistence.GetOneById(ctx, id)
	if err != nil || item.Id == "" {
		return data.BeaconV1{}, err
	}

	if reqctx.OrgId != "" && item.OrgId != reqctx.OrgId {
		return data.BeaconV1{}, nil
	}

	return item, nil
}

func (c *BeaconsMemoryPersistence) Update(ctx context.Context, reqctx cdata.RequestContextV1, item data.BeaconV1) (data.BeaconV1, error) {
	originalItem, err := c.IdentifiableMemoryPersistence.GetOneById(ctx, item.Id)
	if err != nil || originalItem.Id == "" {
		return data.BeaconV1{}, err
	}

	if reqctx.OrgId != "" && originalItem.OrgId != reqctx.OrgId {
		return data.BeaconV1{}, nil
	}

	if reqctx.OrgId != "" {
		item.OrgId = reqctx.OrgId
	}

	return c.IdentifiableMemoryPersistence.Update(ctx, item)
}

func (c *BeaconsMemoryPersistence) DeleteById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.BeaconV1, error) {
	originalItem, err := c.IdentifiableMemoryPersistence.GetOneById(ctx, id)
	if err != nil || originalItem.Id == "" {
		return data.BeaconV1{}, err
	}

	if reqctx.OrgId != "" && originalItem.OrgId != reqctx.OrgId {
		return data.BeaconV1{}, nil
	}

	return c.IdentifiableMemoryPersistence.DeleteById(ctx, id)
}
