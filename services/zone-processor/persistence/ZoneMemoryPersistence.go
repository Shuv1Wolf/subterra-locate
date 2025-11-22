package persistence

import (
	"context"
	"fmt"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type ZoneMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence[data1.ZoneV1, string]
}

func NewZoneMemoryPersistence() *ZoneMemoryPersistence {
	return &ZoneMemoryPersistence{
		IdentifiableMemoryPersistence: *cpersist.NewIdentifiableMemoryPersistence[data1.ZoneV1, string](),
	}
}

func (c *ZoneMemoryPersistence) GetPageByFilter(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[data1.ZoneV1], err error) {
	if reqctx.OrgId != "" {
		filter.Put("org_id", reqctx.OrgId)
	}
	return c.IdentifiableMemoryPersistence.GetPageByFilter(ctx, c.composeFilter(filter), paging, nil, nil)
}

func (c *ZoneMemoryPersistence) GetOneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data1.ZoneV1, error) {
	item, err := c.IdentifiableMemoryPersistence.GetOneById(ctx, id)
	if err != nil || item.Id == "" {
		return item, err
	}
	if item.OrgId != reqctx.OrgId {
		return data1.ZoneV1{}, nil
	}
	return item, nil
}

func (c *ZoneMemoryPersistence) Create(ctx context.Context, reqctx cdata.RequestContextV1, item data1.ZoneV1) (data1.ZoneV1, error) {
	if item.OrgId == "" {
		item.OrgId = reqctx.OrgId
	}
	return c.IdentifiableMemoryPersistence.Create(ctx, item)
}

func (c *ZoneMemoryPersistence) Update(ctx context.Context, reqctx cdata.RequestContextV1, item data1.ZoneV1) (data1.ZoneV1, error) {
	data1, err := c.GetOneById(ctx, reqctx, item.Id)
	if err != nil {
		return data1, err
	}

	if data1.Id == "" {
		return data1, fmt.Errorf("zone not found: %s", item.Id)
	}

	return c.IdentifiableMemoryPersistence.Update(ctx, item)
}

func (c *ZoneMemoryPersistence) DeleteById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data1.ZoneV1, error) {
	data1, err := c.GetOneById(ctx, reqctx, id)
	if err != nil {
		return data1, err
	}

	if data1.Id == "" {
		return data1, fmt.Errorf("zone not found: %s", id)
	}

	return c.IdentifiableMemoryPersistence.DeleteById(ctx, id)
}

func (c *ZoneMemoryPersistence) composeFilter(filter cquery.FilterParams) func(item data1.ZoneV1) bool {
	id, idOk := filter.GetAsNullableString("id")
	orgId, orgIdOk := filter.GetAsNullableString("org_id")
	mapId, mapIdOk := filter.GetAsNullableString("map_id")
	typ, typOk := filter.GetAsNullableString("type")

	return func(item data1.ZoneV1) bool {
		if idOk && item.Id != id {
			return false
		}
		if orgIdOk && item.OrgId != orgId {
			return false
		}
		if mapIdOk && item.MapId != mapId {
			return false
		}
		if typOk && item.Type != typ {
			return false
		}
		return true
	}
}
