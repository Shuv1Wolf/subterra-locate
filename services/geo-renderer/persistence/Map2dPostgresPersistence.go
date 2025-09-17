package persistence

import (
	"context"
	"fmt"
	"strings"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/persistence"
)

type Map2dPostgresPersistence struct {
	cpg.IdentifiablePostgresPersistence[data1.Map2dV1, string]
}

func NewMap2dPostgresPersistence() *Map2dPostgresPersistence {
	c := &Map2dPostgresPersistence{}
	c.IdentifiablePostgresPersistence = *cpg.InheritIdentifiablePostgresPersistence[data1.Map2dV1, string](c, "map-2d")
	c.MaxPageSize = 100
	return c
}

func (c *Map2dPostgresPersistence) DefineSchema() {
	c.ClearSchema()
	c.EnsureSchema("CREATE SEQUENCE IF NOT EXISTS map_2d_id_seq START 100")

	c.EnsureSchema("CREATE TABLE " + c.QuotedTableName() + " (" +
		"\"id\" VARCHAR(32) PRIMARY KEY, " +
		"\"name\" VARCHAR(32), " +
		"\"svg_content\" TEXT, " +
		"\"scale_x\" FLOAT, " +
		"\"scale_y\" FLOAT, " +
		"\"created_at\" TIMESTAMP, " +
		"\"org_id\" VARCHAR(32))")

	c.EnsureIndex(c.TableName+"_org_id", map[string]string{"org_id": "1"}, nil)
}

func (c *Map2dPostgresPersistence) composeFilter(filter cquery.FilterParams) string {
	filters := make([]string, 0)
	if id, ok := filter.GetAsNullableString("id"); ok && id != "" {
		filters = append(filters, "id='"+id+"'")
	}
	if siteId, ok := filter.GetAsNullableString("org_id"); ok && siteId != "" {
		filters = append(filters, "org_id='"+siteId+"'")
	}
	if name, ok := filter.GetAsNullableString("name"); ok && name != "" {
		filters = append(filters, "name='"+name+"'")
	}

	if len(filters) > 0 {
		return strings.Join(filters, " AND ")
	} else {
		return ""
	}
}

func (c *Map2dPostgresPersistence) Create(ctx context.Context, item data1.Map2dV1) (data1.Map2dV1, error) {
	if item.Id == "" {
		var nextId int64
		row := c.Client.QueryRow(ctx, "SELECT nextval('device_id_seq')")
		if err := row.Scan(&nextId); err != nil {
			return item, err
		}
		item.Id = fmt.Sprintf("map2d$%d", nextId)
	}

	return c.IdentifiablePostgresPersistence.Create(ctx, item)
}

func (c *Map2dPostgresPersistence) GetPageByFilter(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data1.Map2dV1], error) {
	return c.IdentifiablePostgresPersistence.GetPageByFilter(ctx,
		c.composeFilter(filter), paging,
		"", "",
	)
}
