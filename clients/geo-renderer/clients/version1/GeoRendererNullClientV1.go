package clients1

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type GeoRendererNullClientV1 struct {
}

func NewGeoRendererNullClientV1() *GeoRendererNullClientV1 {
	return &GeoRendererNullClientV1{}
}

func (c *GeoRendererNullClientV1) GetMaps(ctx context.Context, reqctx cdata.RequestContextV1,
	filter *cquery.FilterParams,
	paging *cquery.PagingParams) (*cquery.DataPage[data1.Map2dV1], error) {
	return cquery.NewEmptyDataPage[data1.Map2dV1](), nil
}

func (c *GeoRendererNullClientV1) GetMapById(ctx context.Context, reqctx cdata.RequestContextV1,
	id string) (*data1.Map2dV1, error) {
	return nil, nil
}

func (c *GeoRendererNullClientV1) CreateMap(ctx context.Context, reqctx cdata.RequestContextV1,
	map2d data1.Map2dV1) (*data1.Map2dV1, error) {
	return nil, nil
}

func (c *GeoRendererNullClientV1) UpdateMap(ctx context.Context, reqctx cdata.RequestContextV1,
	map2d data1.Map2dV1) (*data1.Map2dV1, error) {
	return nil, nil
}

func (c *GeoRendererNullClientV1) DeleteMapById(ctx context.Context, reqctx cdata.RequestContextV1,
	id string) (*data1.Map2dV1, error) {
	return nil, nil
}

func (c *GeoRendererNullClientV1) GetZones(ctx context.Context, reqctx cdata.RequestContextV1,
	filter *cquery.FilterParams,
	paging *cquery.PagingParams) (*cquery.DataPage[data1.ZoneV1], error) {
	return cquery.NewEmptyDataPage[data1.ZoneV1](), nil
}

func (c *GeoRendererNullClientV1) GetZoneById(ctx context.Context, reqctx cdata.RequestContextV1,
	id string) (*data1.ZoneV1, error) {
	return nil, nil
}

func (c *GeoRendererNullClientV1) CreateZone(ctx context.Context, reqctx cdata.RequestContextV1,
	zone data1.ZoneV1) (*data1.ZoneV1, error) {
	return nil, nil
}

func (c *GeoRendererNullClientV1) UpdateZone(ctx context.Context, reqctx cdata.RequestContextV1,
	zone data1.ZoneV1) (*data1.ZoneV1, error) {
	return nil, nil
}

func (c *GeoRendererNullClientV1) DeleteZoneById(ctx context.Context, reqctx cdata.RequestContextV1,
	id string) (*data1.ZoneV1, error) {
	return nil, nil
}
