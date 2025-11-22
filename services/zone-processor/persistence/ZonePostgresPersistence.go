package persistence

import (
	"context"
	"fmt"
	"strings"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/persistence"
)

type ZonePostgresPersistence struct {
	cpg.IdentifiablePostgresPersistence[data1.ZoneV1, string]
}

func NewZonePostgresPersistence() *ZonePostgresPersistence {
	c := &ZonePostgresPersistence{}
	c.IdentifiablePostgresPersistence = *cpg.InheritIdentifiablePostgresPersistence[data1.ZoneV1, string](c, "zone")
	c.MaxPageSize = 100
	return c
}

func (c *ZonePostgresPersistence) DefineSchema() {
	c.ClearSchema()
	c.EnsureSchema("CREATE SEQUENCE IF NOT EXISTS zone_id_seq START 100")

	c.EnsureSchema("CREATE TABLE " + c.QuotedTableName() + " (" +
		"\"id\" VARCHAR(32) PRIMARY KEY, " +
		"\"map_id\" VARCHAR(32), " +
		"\"org_id\" VARCHAR(32), " +
		"\"name\" VARCHAR(32), " +
		"\"position_x\" FLOAT, " +
		"\"position_y\" FLOAT, " +
		"\"width\" FLOAT, " +
		"\"height\" FLOAT, " +
		"\"max_device\" INTEGER, " +
		"\"type\" VARCHAR(32), " +
		"\"created_at\" TIMESTAMP)")

	c.EnsureIndex(c.TableName+"_org_id", map[string]string{"org_id": "1"}, nil)
	c.EnsureIndex(c.TableName+"_map_id", map[string]string{"map_id": "1"}, nil)
}

func (c *ZonePostgresPersistence) composeFilter(filter cquery.FilterParams) string {
	filters := make([]string, 0)
	if id, ok := filter.GetAsNullableString("id"); ok && id != "" {
		filters = append(filters, "id='"+id+"'")
	}
	if orgId, ok := filter.GetAsNullableString("org_id"); ok && orgId != "" {
		filters = append(filters, "org_id='"+orgId+"'")
	}
	if mapId, ok := filter.GetAsNullableString("map_id"); ok && mapId != "" {
		filters = append(filters, "map_id='"+mapId+"'")
	}
	if name, ok := filter.GetAsNullableString("name"); ok && name != "" {
		filters = append(filters, "name='"+name+"'")
	}
	if typ, ok := filter.GetAsNullableString("type"); ok && typ != "" {
		filters = append(filters, "type='"+typ+"'")
	}

	if len(filters) > 0 {
		return strings.Join(filters, " AND ")
	} else {
		return ""
	}
}

func (c *ZonePostgresPersistence) Create(ctx context.Context, reqctx cdata.RequestContextV1, item data1.ZoneV1) (data1.ZoneV1, error) {
	if item.Id == "" {
		var nextId int64
		row := c.Client.QueryRow(ctx, "SELECT nextval('zone_id_seq')")
		if err := row.Scan(&nextId); err != nil {
			return item, err
		}
		item.Id = fmt.Sprintf("zone$%d", nextId)
	}

	if item.OrgId == "" {
		item.OrgId = reqctx.OrgId
	}

	return c.IdentifiablePostgresPersistence.Create(ctx, item)
}

func (c *ZonePostgresPersistence) GetPageByFilter(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data1.ZoneV1], error) {
	if reqctx.OrgId != "" {
		filter.Put("org_id", reqctx.OrgId)
	}
	return c.IdentifiablePostgresPersistence.GetPageByFilter(ctx,
		c.composeFilter(filter), paging,
		"", "",
	)
}

func (c *ZonePostgresPersistence) GetOneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data1.ZoneV1, error) {
	filter := cquery.NewFilterParamsFromTuples(
		"id", id,
	)

	if reqctx.OrgId != "" {
		filter.Put("org_id", reqctx.OrgId)
	}

	paging := *cquery.NewPagingParams(0, 1, false)
	page, err := c.IdentifiablePostgresPersistence.GetPageByFilter(ctx,
		c.composeFilter(*filter), paging,
		"", "",
	)
	if err != nil {
		return data1.ZoneV1{}, err
	}
	if page.HasData() {
		return page.Data[0], nil
	}
	return data1.ZoneV1{}, nil
}

func (c *ZonePostgresPersistence) Update(ctx context.Context, reqctx cdata.RequestContextV1, item data1.ZoneV1) (data1.ZoneV1, error) {
	data, err := c.GetOneById(ctx, reqctx, item.Id)
	if err != nil {
		return data, err
	}

	if data.Id == "" {
		return data, fmt.Errorf("zone not found: %s", item.Id)
	}

	return c.IdentifiablePostgresPersistence.Update(ctx, item)
}

func (c *ZonePostgresPersistence) DeleteById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data1.ZoneV1, error) {
	data, err := c.GetOneById(ctx, reqctx, id)
	if err != nil {
		return data, err
	}

	if data.Id == "" {
		return data, fmt.Errorf("zone not found: %s", id)
	}

	return c.IdentifiablePostgresPersistence.DeleteById(ctx, id)
}
