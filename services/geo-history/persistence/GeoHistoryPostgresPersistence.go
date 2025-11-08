package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-history/data/version1"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	ccdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/persistence"
	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/persistence"
)

type GeoHistoryPostgresPersistence struct {
	cpg.IdentifiablePostgresPersistence[data1.HistoricalRecordV1, string]
	retentionPolicyInterval string
	chunkSize               string
	compressionInterval     string
	retentionPolicyEnabled  bool

	isOpen bool
}

func NewGeoHistoryPostgresPersistence() (*GeoHistoryPostgresPersistence, error) {
	c := &GeoHistoryPostgresPersistence{}
	c.IdentifiablePostgresPersistence = *cpg.InheritIdentifiablePostgresPersistence[data1.HistoricalRecordV1, string](c, "geo-history")
	c.MaxPageSize = 100

	c.retentionPolicyInterval = "30 weeks"
	c.chunkSize = "1 days"
	c.compressionInterval = "4 weeks"
	c.retentionPolicyEnabled = true

	return c, nil
}

func (c *GeoHistoryPostgresPersistence) Open(ctx context.Context) error {
	if c.isOpen {
		return nil
	}

	if err := c.IdentifiablePostgresPersistence.Open(ctx); err != nil {
		return err
	}

	// Timescale init
	if err := c.addCustomIndex(ctx); err != nil {
		return err
	}
	if err := c.createHypertable(ctx, c.chunkSize); err != nil {
		return err
	}
	if err := c.addCompression(ctx, c.compressionInterval); err != nil {
		return err
	}
	if c.retentionPolicyEnabled {
		if err := c.addRetentionPolicy(ctx, c.retentionPolicyInterval); err != nil {
			return err
		}
	}
	c.isOpen = true
	return nil
}

func (c *GeoHistoryPostgresPersistence) IsOpen() bool {
	return c.isOpen
}

func (c *GeoHistoryPostgresPersistence) Close(ctx context.Context) error {
	if c.isOpen {
		c.PostgresPersistence.Close(ctx)
		c.isOpen = false
	}
	return nil
}

func (c *GeoHistoryPostgresPersistence) DefineSchema() {
	c.ClearSchema()

	c.EnsureSchema("CREATE TABLE IF NOT EXISTS " + c.QuotedTableName() + " (" +
		"\"id\" VARCHAR(32), " +
		"\"entity_id\" VARCHAR(50), " +
		"\"timestamp\" TIMESTAMP NOT NULL, " +
		"\"x\" DOUBLE PRECISION, " +
		"\"y\" DOUBLE PRECISION, " +
		"\"z\" DOUBLE PRECISION, " +
		"\"org_id\" VARCHAR(32), " +
		"\"map_id\" VARCHAR(32)" +
		")")

	c.EnsureIndex(c.TableName+"_tm_org_id_map_id", map[string]string{"timestamp": "1", "org_id": "1", "map_id": "1"}, nil)
	c.EnsureIndex(c.TableName+"_tm_entity_id", map[string]string{"timestamp": "1", "entity_id": "1"}, nil)
}

func (c *GeoHistoryPostgresPersistence) addRetentionPolicy(ctx context.Context, interval string) error {
	schema := c.SchemaName
	if schema == "" {
		schema = "public"
	}

	fullName := `"` + schema + `"."` + c.TableName + `"`
	query := "SELECT add_retention_policy('" + fullName + "', INTERVAL '" + interval + "', if_not_exists => TRUE)"

	result, err := c.Client.Query(ctx, query)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to create retention_policy")
		return err
	}
	defer result.Close()

	return result.Err()
}

func (c *GeoHistoryPostgresPersistence) addCustomIndex(ctx context.Context) error {
	query := "CREATE INDEX IF NOT EXISTS \"" + c.TableName + "_tm_org_id_map_id\" ON " +
		c.QuotedTableName() + `(timestamp, org_id, map_id) 
		WHERE org_id IS NOT NULL AND map_id IS NOT NULL`

	result, err := c.Client.Query(ctx, query)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to create custom index")
		return err
	}
	result.Close()

	query = "CREATE INDEX IF NOT EXISTS \"" + c.TableName + "_tm_entity_id\" ON " +
		c.QuotedTableName() + `(timestamp, entity_id) 
		WHERE entity_id IS NOT NULL`

	result, err = c.Client.Query(ctx, query)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to create custom index")
		return err
	}
	result.Close()

	return result.Err()
}

func (c *GeoHistoryPostgresPersistence) createHypertable(ctx context.Context, chunkInterval string) error {
	schema := c.SchemaName
	if schema == "" {
		schema = "public"
	}
	fullName := `"` + schema + `"."` + c.TableName + `"`

	// Проверим, есть ли уже hypertable
	check := "SELECT hypertable_name FROM timescaledb_information.hypertables WHERE hypertable_schema = '" +
		schema + "' AND hypertable_name = '" + c.TableName + "'"
	res, err := c.Client.Query(ctx, check)
	if err != nil {
		return err
	}
	if res.Next() {
		c.Logger.Info(ctx, "Hypertable already exists: %s", fullName)
		res.Close()
		return nil
	}
	res.Close()

	// Создаём hypertable
	query := "SELECT create_hypertable('" + fullName +
		"', 'timestamp', chunk_time_interval => INTERVAL '" + chunkInterval +
		"', migrate_data => TRUE, if_not_exists => TRUE)"
	result, err := c.Client.Query(ctx, query)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to autocreate hypertable")
		return err
	}
	result.Close()

	// Проверим, что Timescale поддерживает compression параметры
	enableCompression := "ALTER TABLE " + fullName + ` SET (
		timescaledb.compress,
		timescaledb.compress_orderby = 'timestamp',
		timescaledb.compress_segmentby = 'entity_id'
	);`
	result, err = c.Client.Query(ctx, enableCompression)
	if err != nil {
		c.Logger.Warn(ctx, "Compression setup skipped: %v", err)
		return nil
	}
	result.Close()

	c.Logger.Info(ctx, "Hypertable created with compression: %s", fullName)
	return nil
}

func (c *GeoHistoryPostgresPersistence) addCompression(ctx context.Context, compressionInterval string) error {
	schema := c.SchemaName
	if schema == "" {
		schema = "public"
	}

	fullName := `"` + schema + `"."` + c.TableName + `"`
	query := "SELECT add_compression_policy('" + fullName + "', INTERVAL '" + compressionInterval + "', if_not_exists => TRUE)"

	result, err := c.Client.Query(ctx, query)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to set compression policy")
		return err
	}
	defer result.Close()

	return result.Err()
}

func (c *GeoHistoryPostgresPersistence) composeFilter(filter *cquery.FilterParams) string {
	if filter == nil {
		return ""
	}

	filters := make([]string, 0)

	if to, ok := filter.GetAsNullableInteger("to"); ok {
		filters = append(filters, fmt.Sprintf("(EXTRACT(EPOCH FROM (timestamp)))::int <= '%d'", to))
	}
	if from, ok := filter.GetAsNullableInteger("from"); ok {
		filters = append(filters, fmt.Sprintf("(EXTRACT(EPOCH FROM (timestamp)))::int >= '%d'", from))
	}
	if orgId, ok := filter.GetAsNullableString("org_id"); ok && orgId != "" {
		filters = append(filters, "org_id='"+orgId+"'")
	}
	if mapId, ok := filter.GetAsNullableString("map_id"); ok && mapId != "" {
		filters = append(filters, "map_id='"+mapId+"'")
	}
	if entityId, ok := filter.GetAsNullableString("entity_id"); ok && entityId != "" {
		filters = append(filters, "entity_id='"+entityId+"'")
	}

	var entityIds []string
	if _idis, ok := filter.GetAsObject("entity_ids"); ok {
		if _val, ok := _idis.([]string); ok {
			entityIds = _val
		}
		if _val, ok := _idis.(string); ok {
			entityIds = strings.Split(_val, ",")
		}
		if entityIds[0] != "" && len(entityIds) != 0 {
			filters = append(filters, "entity_id IN ('"+strings.Join(entityIds, "','")+"')")
		}

	}

	if len(filters) > 0 {
		return strings.Join(filters, " AND ")
	} else {
		return ""
	}
}

func (c *GeoHistoryPostgresPersistence) InsertBatch(ctx context.Context, items []*data1.HistoricalRecordV1) error {
	batch := &pgx.Batch{}

	for _, item := range items {
		newItem := c.cloneItem(*item)
		newItem = persistence.GenerateObjectIdIfNotExists[*data1.HistoricalRecordV1](newItem)

		objMap, convErr := c.Overrides.ConvertFromPublic(*newItem)
		if convErr != nil {
			return convErr
		}
		columns, values := c.GenerateColumnsAndValues(objMap)

		columnsStr := c.GenerateColumns(columns)
		paramsStr := c.GenerateParameters(len(values))

		query := "INSERT INTO " + c.QuotedTableName() +
			" (" + columnsStr + ") VALUES (" + paramsStr + ") "
		batch.Queue(query, values...)
	}

	results := c.Client.SendBatch(ctx, batch)
	defer results.Close()

	for _, item := range items {
		_, err := results.Exec()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				c.Logger.Error(ctx, err, "Item %s already exists", item.Id)
				continue
			}

			return fmt.Errorf("unable to insert row: %v", err)
		}
	}

	return results.Close()
}

func (c *GeoHistoryPostgresPersistence) GetHistory(ctx context.Context, reqctx ccdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams, sortField cquery.SortField) (cquery.DataPage[data1.HistoricalRecordV1], error) {
	sort := ""
	if sortField.Name != "" {
		sort += sortField.Name
	}

	if sortField.Ascending {
		sort += " ASC"
	} else {
		sort += " DESC"
	}

	if reqctx.OrgId != "" {
		filter.Put("org_id", reqctx.OrgId)
	}

	return c.IdentifiablePostgresPersistence.GetPageByFilter(ctx,
		c.composeFilter(&filter), paging, sort, "",
	)
}

func (c *GeoHistoryPostgresPersistence) cloneItem(item any) *data1.HistoricalRecordV1 {
	if cloneableItem, ok := item.(cdata.ICloneable[*data1.HistoricalRecordV1]); ok {
		return cloneableItem.Clone()
	}

	strObject, _ := c.JsonConvertor.ToJson(item.(data1.HistoricalRecordV1))
	newItem, _ := c.JsonConvertor.FromJson(strObject)
	return &newItem
}
