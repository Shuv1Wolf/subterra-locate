package clients1

import (
	"context"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
	ccdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
)

type ZoneGrpcClientV1 struct {
	*cclients.CommandableGrpcClient
}

func NewZoneGrpcClientV1() *ZoneGrpcClientV1 {
	c := &ZoneGrpcClientV1{
		CommandableGrpcClient: cclients.NewCommandableGrpcClient("zone_processor.v1"),
	}
	return c
}

func (c *ZoneGrpcClientV1) GetZones(ctx context.Context, reqctx cdata.RequestContextV1, filter *cquery.FilterParams, paging *cquery.PagingParams) (*cquery.DataPage[data1.ZoneV1], error) {
	var pagingMap map[string]interface{}
	var reqctxMap map[string]interface{}

	if paging != nil {
		pagingMap = map[string]interface{}{
			"skip":  paging.Skip,
			"take":  paging.Take,
			"total": paging.Total,
		}
	}

	reqctxMap = map[string]interface{}{
		"org_id":  reqctx.OrgId,
		"user_id": reqctx.UserId,
	}
	params := ccdata.NewAnyValueMapFromTuples(
		"filter", filter.Value(),
		"paging", pagingMap,
		"reqctx", reqctxMap,
	)
	response, err := c.CallCommand(ctx, "get_zones", ccdata.NewAnyValueMapFromValue(params.Value()))
	if err != nil {
		return cquery.NewEmptyDataPage[data1.ZoneV1](), err
	}
	return cclients.HandleHttpResponse[*cquery.DataPage[data1.ZoneV1]](response)
}

func (c *ZoneGrpcClientV1) GetZoneById(ctx context.Context, reqctx cdata.RequestContextV1, zoneId string) (*data1.ZoneV1, error) {
	var reqctxMap map[string]interface{}
	reqctxMap = map[string]interface{}{
		"org_id":  reqctx.OrgId,
		"user_id": reqctx.UserId,
	}

	params := ccdata.NewAnyValueMapFromTuples(
		"zone_id", zoneId,
		"reqctx", reqctxMap,
	)

	response, err := c.CallCommand(ctx, "get_zone_by_id", params)
	if err != nil {
		return nil, err
	}
	return cclients.HandleHttpResponse[*data1.ZoneV1](response)
}

func (c *ZoneGrpcClientV1) CreateZone(ctx context.Context, reqctx cdata.RequestContextV1, zone data1.ZoneV1) (*data1.ZoneV1, error) {
	var reqctxMap map[string]interface{}
	reqctxMap = map[string]interface{}{
		"org_id":  reqctx.OrgId,
		"user_id": reqctx.UserId,
	}
	params := ccdata.NewAnyValueMapFromTuples(
		"zone", zone,
		"reqctx", reqctxMap,
	)
	response, err := c.CallCommand(ctx, "create_zone", params)
	if err != nil {
		return nil, err
	}
	return cclients.HandleHttpResponse[*data1.ZoneV1](response)
}

func (c *ZoneGrpcClientV1) UpdateZone(ctx context.Context, reqctx cdata.RequestContextV1, zone data1.ZoneV1) (*data1.ZoneV1, error) {
	var reqctxMap map[string]interface{}
	reqctxMap = map[string]interface{}{
		"org_id":  reqctx.OrgId,
		"user_id": reqctx.UserId,
	}
	params := ccdata.NewAnyValueMapFromTuples(
		"zone", zone,
		"reqctx", reqctxMap,
	)
	response, err := c.CallCommand(ctx, "update_zone", params)
	if err != nil {
		return nil, err
	}
	return cclients.HandleHttpResponse[*data1.ZoneV1](response)
}

func (c *ZoneGrpcClientV1) DeleteZoneById(ctx context.Context, reqctx cdata.RequestContextV1, zoneId string) (*data1.ZoneV1, error) {
	var reqctxMap map[string]interface{}
	reqctxMap = map[string]interface{}{
		"org_id":  reqctx.OrgId,
		"user_id": reqctx.UserId,
	}
	params := ccdata.NewAnyValueMapFromTuples(
		"zone_id", zoneId,
		"reqctx", reqctxMap,
	)
	response, err := c.CallCommand(ctx, "delete_zone_by_id", params)
	if err != nil {
		return nil, err
	}
	return cclients.HandleHttpResponse[*data1.ZoneV1](response)
}
