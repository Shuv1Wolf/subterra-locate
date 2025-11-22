package logic

import (
	"context"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IMapService interface {

	// MAP
	GetMaps(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data1.Map2dV1], error)

	GetMapById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data1.Map2dV1, error)

	CreateMap(ctx context.Context, reqctx cdata.RequestContextV1, map2d data1.Map2dV1) (data1.Map2dV1, error)

	UpdateMap(ctx context.Context, reqctx cdata.RequestContextV1, map2d data1.Map2dV1) (data1.Map2dV1, error)

	DeleteMapById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data1.Map2dV1, error)

	// ZONE

	GetZones(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data1.ZoneV1], error)

	GetZoneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data1.ZoneV1, error)

	CreateZone(ctx context.Context, reqctx cdata.RequestContextV1, zone data1.ZoneV1) (data1.ZoneV1, error)

	UpdateZone(ctx context.Context, reqctx cdata.RequestContextV1, zone data1.ZoneV1) (data1.ZoneV1, error)

	DeleteZoneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data1.ZoneV1, error)
}
