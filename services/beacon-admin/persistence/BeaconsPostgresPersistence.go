package persistence

import (
	"context"
	"fmt"
	"strings"

	data "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/persistence"
)

type BeaconsPostgresPersistence struct {
	cpg.IdentifiablePostgresPersistence[data.BeaconV1, string]
}

func NewBeaconsPostgresPersistence() *BeaconsPostgresPersistence {
	c := &BeaconsPostgresPersistence{}
	c.IdentifiablePostgresPersistence = *cpg.InheritIdentifiablePostgresPersistence[data.BeaconV1, string](c, "beacon-admin")
	c.MaxPageSize = 100
	return c
}

func (c *BeaconsPostgresPersistence) DefineSchema() {
	c.ClearSchema()
	c.EnsureSchema("CREATE SEQUENCE IF NOT EXISTS beacon_id_seq START 100")

	c.EnsureSchema("CREATE TABLE " + c.QuotedTableName() + " (" +
		"\"id\" VARCHAR(32) PRIMARY KEY, " +
		"\"type\" VARCHAR(15), " +
		"\"udi\" VARCHAR(50), " +
		"\"label\" VARCHAR(50), " +
		"\"x\" FLOAT, " +
		"\"y\" FLOAT, " +
		"\"z\" FLOAT, " +
		"\"org_id\" VARCHAR(32), " +
		"\"map_id\" VARCHAR(32), " +
		"\"enabled\" BOOLEAN)")

	c.EnsureIndex(c.TableName+"_type", map[string]string{"type": "1"}, nil)
	c.EnsureIndex(c.TableName+"_org_id", map[string]string{"org_id": "1"}, nil)
	c.EnsureIndex(c.TableName+"_map_id", map[string]string{"map_id": "1"}, nil)
	c.EnsureIndex(c.TableName+"_udi", map[string]string{"udi": "1"}, nil)
}

func (c *BeaconsPostgresPersistence) composeFilter(filter cquery.FilterParams) string {
	filters := make([]string, 0)
	if id, ok := filter.GetAsNullableString("id"); ok && id != "" {
		filters = append(filters, "id='"+id+"'")
	}
	if org_id, ok := filter.GetAsNullableString("org_id"); ok && org_id != "" {
		filters = append(filters, "org_id='"+org_id+"'")
	}
	if map_id, ok := filter.GetAsNullableString("map_id"); ok && map_id != "" {
		filters = append(filters, "map_id='"+map_id+"'")
	}
	if typeId, ok := filter.GetAsNullableString("type"); ok && typeId != "" {
		filters = append(filters, "type='"+typeId+"'")
	}
	if udi, ok := filter.GetAsNullableString("udi"); ok && udi != "" {
		filters = append(filters, "udi='"+udi+"'")
	}
	if label, ok := filter.GetAsNullableString("label"); ok && label != "" {
		filters = append(filters, "label='"+label+"'")
	}
	if udis, ok := filter.GetAsNullableString("udis"); ok {
		ids := strings.Split(udis, ",")
		filters = append(filters, "udi IN ('"+strings.Join(ids, "','")+"')")
	}
	if enabled, ok := filter.GetAsNullableString("enabled"); ok && enabled != "" {
		filters = append(filters, "enabled="+enabled)
	}

	if len(filters) > 0 {
		return strings.Join(filters, " AND ")
	} else {
		return ""
	}
}

func (c *BeaconsPostgresPersistence) Create(ctx context.Context, reqctx cdata.RequestContextV1, item data.BeaconV1) (data.BeaconV1, error) {
	if item.Id == "" {
		var nextId int64
		row := c.Client.QueryRow(ctx, "SELECT nextval('beacon_id_seq')")
		if err := row.Scan(&nextId); err != nil {
			return item, err
		}
		item.Id = fmt.Sprintf("beacon$%d", nextId)
	}

	if item.OrgId != "" {
		item.OrgId = reqctx.OrgId
	}

	return c.IdentifiablePostgresPersistence.Create(ctx, item)
}

func (c *BeaconsPostgresPersistence) GetPageByFilter(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.BeaconV1], error) {
	if reqctx.OrgId != "" {
		filter.Put("org_id", reqctx.OrgId)
	}

	return c.IdentifiablePostgresPersistence.GetPageByFilter(ctx,
		c.composeFilter(filter), paging,
		"", "",
	)
}

func (c *BeaconsPostgresPersistence) GetOneByUdi(ctx context.Context, reqctx cdata.RequestContextV1, udi string) (data.BeaconV1, error) {
	filter := cquery.NewFilterParamsFromTuples(
		"udi", udi,
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
		return data.BeaconV1{}, err
	}
	if page.HasData() {
		return page.Data[0], nil
	}
	return data.BeaconV1{}, nil
}

func (c *BeaconsPostgresPersistence) GetOneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.BeaconV1, error) {
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
		return data.BeaconV1{}, err
	}
	if page.HasData() {
		return page.Data[0], nil
	}
	return data.BeaconV1{}, nil
}

func (c *BeaconsPostgresPersistence) Update(ctx context.Context, reqctx cdata.RequestContextV1, item data.BeaconV1) (data.BeaconV1, error) {
	data, err := c.GetOneById(ctx, reqctx, item.Id)
	if err != nil {
		return data, err
	}

	if data.Id != "" {
		return data, fmt.Errorf("beacon not found: %s", item.Id)
	}

	return c.IdentifiablePostgresPersistence.Update(ctx, item)
}

func (c *BeaconsPostgresPersistence) DeleteById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.BeaconV1, error) {
	data, err := c.GetOneById(ctx, reqctx, id)
	if err != nil {
		return data, err
	}

	if data.Id != "" {
		return data, fmt.Errorf("beacon not found: %s", id)
	}

	return c.IdentifiablePostgresPersistence.DeleteById(ctx, id)
}
