package persistence

import (
	"context"
	"fmt"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type Map2dMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence[data.Map2dV1, string]
}

func NewMap2dMemoryPersistence() *Map2dMemoryPersistence {
	return &Map2dMemoryPersistence{
		IdentifiableMemoryPersistence: *cpersist.NewIdentifiableMemoryPersistence[data.Map2dV1, string](),
	}
}

func (c *Map2dMemoryPersistence) GetPageByFilter(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[data.Map2dV1], err error) {
	if reqctx.OrgId != "" {
		filter.Put("org_id", reqctx.OrgId)
	}
	return c.IdentifiableMemoryPersistence.GetPageByFilter(ctx, c.composeFilter(filter), paging, nil, nil)
}

func (c *Map2dMemoryPersistence) GetOneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.Map2dV1, error) {
	item, err := c.IdentifiableMemoryPersistence.GetOneById(ctx, id)
	if err != nil || item.Id == "" {
		return item, err
	}
	if item.OrgId != reqctx.OrgId {
		return data.Map2dV1{}, nil
	}
	return item, nil
}

func (c *Map2dMemoryPersistence) Create(ctx context.Context, reqctx cdata.RequestContextV1, item data.Map2dV1) (data.Map2dV1, error) {
	if item.OrgId == "" {
		item.OrgId = reqctx.OrgId
	}
	return c.IdentifiableMemoryPersistence.Create(ctx, item)
}

func (c *Map2dMemoryPersistence) Update(ctx context.Context, reqctx cdata.RequestContextV1, item data.Map2dV1) (data.Map2dV1, error) {
	data, err := c.GetOneById(ctx, reqctx, item.Id)
	if err != nil {
		return data, err
	}

	if data.Id == "" {
		return data, fmt.Errorf("map not found: %s", item.Id)
	}

	return c.IdentifiableMemoryPersistence.Update(ctx, item)
}

func (c *Map2dMemoryPersistence) DeleteById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.Map2dV1, error) {
	data, err := c.GetOneById(ctx, reqctx, id)
	if err != nil {
		return data, err
	}

	if data.Id == "" {
		return data, fmt.Errorf("map not found: %s", id)
	}

	return c.IdentifiableMemoryPersistence.DeleteById(ctx, id)
}

func (c *Map2dMemoryPersistence) composeFilter(filter cquery.FilterParams) func(item data.Map2dV1) bool {
	id, idOk := filter.GetAsNullableString("id")
	orgId, orgIdOk := filter.GetAsNullableString("org_id")

	return func(item data.Map2dV1) bool {
		if idOk && item.Id != id {
			return false
		}
		if orgIdOk && item.OrgId != orgId {
			return false
		}
		return true
	}
}
