package persistence

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type Map2dMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence[data.Map2dV1, string]
}

func NewMap2dMemoryPersistence() *Map2dMemoryPersistence {
	return &Map2dMemoryPersistence{
		IdentifiableMemoryPersistence: *cpersist.NewIdentifiableMemoryPersistence[data.Map2dV1, string](),
	}
}

func (c *Map2dMemoryPersistence) GetPageByFilter(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[data.Map2dV1], err error) {
	
	return c.IdentifiableMemoryPersistence.GetPageByFilter(ctx, c.composeFilter(filter), paging, nil, nil)
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
