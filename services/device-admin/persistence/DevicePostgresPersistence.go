package persistence

import (
	"context"
	"fmt"
	"strings"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/persistence"
)

type DevicePostgresPersistence struct {
	cpg.IdentifiablePostgresPersistence[data.DeviceV1, string]
}

func NewDevicePostgresPersistence() *DevicePostgresPersistence {
	c := &DevicePostgresPersistence{}
	c.IdentifiablePostgresPersistence = *cpg.InheritIdentifiablePostgresPersistence[data.DeviceV1, string](c, "device-admin")
	c.MaxPageSize = 100
	return c
}

func (c *DevicePostgresPersistence) DefineSchema() {
	c.ClearSchema()
	c.EnsureSchema("CREATE SEQUENCE IF NOT EXISTS device_id_seq START 100")

	c.EnsureSchema("CREATE TABLE " + c.QuotedTableName() + " (" +
		"\"id\" VARCHAR(32) PRIMARY KEY, " +
		"\"type\" VARCHAR(15), " +
		"\"name\" VARCHAR(50), " +
		"\"model\" VARCHAR(50), " +
		"\"org_id\" VARCHAR(32), " +
		"\"mac_address\" VARCHAR(50), " +
		"\"ip_address\" VARCHAR(50), " +
		"\"enabled\" BOOLEAN)")

	c.EnsureIndex(c.TableName+"_type", map[string]string{"type": "1"}, nil)
	c.EnsureIndex(c.TableName+"_org_id", map[string]string{"org_id": "1"}, nil)
	c.EnsureIndex(c.TableName+"_model", map[string]string{"model": "1"}, nil)
}

func (c *DevicePostgresPersistence) composeFilter(filter cquery.FilterParams) string {
	filters := make([]string, 0)
	if id, ok := filter.GetAsNullableString("id"); ok && id != "" {
		filters = append(filters, "id='"+id+"'")
	}
	if siteId, ok := filter.GetAsNullableString("org_id"); ok && siteId != "" {
		filters = append(filters, "org_id='"+siteId+"'")
	}
	if typeId, ok := filter.GetAsNullableString("type"); ok && typeId != "" {
		filters = append(filters, "type='"+typeId+"'")
	}
	if name, ok := filter.GetAsNullableString("name"); ok && name != "" {
		filters = append(filters, "name='"+name+"'")
	}
	if model, ok := filter.GetAsNullableString("model"); ok && model != "" {
		filters = append(filters, "model='"+model+"'")
	}
	if enabled, ok := filter.GetAsNullableString("enabled"); ok && enabled != "" {
		filters = append(filters, "enabled="+enabled)
	}
	if macAddress, ok := filter.GetAsNullableString("mac_address"); ok && macAddress != "" {
		filters = append(filters, "mac_address='"+macAddress+"'")
	}

	if len(filters) > 0 {
		return strings.Join(filters, " AND ")
	} else {
		return ""
	}
}

func (c *DevicePostgresPersistence) Create(ctx context.Context, reqctx cdata.RequestContextV1, item data.DeviceV1) (data.DeviceV1, error) {
	if item.Id == "" {
		var nextId int64
		row := c.Client.QueryRow(ctx, "SELECT nextval('device_id_seq')")
		if err := row.Scan(&nextId); err != nil {
			return item, err
		}
		item.Id = fmt.Sprintf("device$%d", nextId)
	}

	if item.OrgId != "" {
		item.OrgId = reqctx.OrgId
	}

	return c.IdentifiablePostgresPersistence.Create(ctx, item)
}

func (c *DevicePostgresPersistence) GetPageByFilter(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.DeviceV1], error) {
	if reqctx.OrgId != "" {
		filter.Put("org_id", reqctx.OrgId)
	}

	return c.IdentifiablePostgresPersistence.GetPageByFilter(ctx,
		c.composeFilter(filter), paging,
		"", "",
	)
}

func (c *DevicePostgresPersistence) GetOneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.DeviceV1, error) {
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
		return data.DeviceV1{}, err
	}
	if page.HasData() {
		return page.Data[0], nil
	}
	return data.DeviceV1{}, nil
}

func (c *DevicePostgresPersistence) Update(ctx context.Context, reqctx cdata.RequestContextV1, item data.DeviceV1) (data.DeviceV1, error) {
	data, err := c.GetOneById(ctx, reqctx, item.Id)
	if err != nil {
		return data, err
	}

	if data.Id != "" {
		return data, fmt.Errorf("device not found: %s", item.Id)
	}

	return c.IdentifiablePostgresPersistence.Update(ctx, item)
}

func (c *DevicePostgresPersistence) DeleteById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.DeviceV1, error) {
	data, err := c.GetOneById(ctx, reqctx, id)
	if err != nil {
		return data, err
	}

	if data.Id != "" {
		return data, fmt.Errorf("device not found: %s", id)
	}

	return c.IdentifiablePostgresPersistence.DeleteById(ctx, id)
}
