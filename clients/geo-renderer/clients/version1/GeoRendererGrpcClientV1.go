package clients1

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
)

type GeoRendererGrpcClientV1 struct {
	*cclients.CommandableGrpcClient
}

func NewGeoRendererGrpcClientV1() *GeoRendererGrpcClientV1 {
	c := &GeoRendererGrpcClientV1{
		CommandableGrpcClient: cclients.NewCommandableGrpcClient("geo.renderer.v1"),
	}
	return c
}

func (c *GeoRendererGrpcClientV1) GetMaps(ctx context.Context,
	filter *cquery.FilterParams,
	paging *cquery.PagingParams) (*cquery.DataPage[data1.Map2dV1], error) {

	var pagingMap map[string]interface{}

	if paging != nil {
		pagingMap = map[string]interface{}{
			"skip":  paging.Skip,
			"take":  paging.Take,
			"total": paging.Total,
		}
	}

	params := cdata.NewAnyValueMapFromTuples(
		"filter", filter.StringValueMap.Value(),
		"paging", pagingMap,
	)

	response, err := c.CallCommand(ctx, "get_maps", cdata.NewAnyValueMapFromValue(params.Value()))

	if err != nil {
		return cquery.NewEmptyDataPage[data1.Map2dV1](), err
	}

	return cclients.HandleHttpResponse[*cquery.DataPage[data1.Map2dV1]](response)
}

func (c *GeoRendererGrpcClientV1) GetMapById(ctx context.Context,
	id string) (*data1.Map2dV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"map_id", id,
	)

	response, err := c.CallCommand(ctx, "get_map_by_id", params)

	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.Map2dV1](response)
}

func (c *GeoRendererGrpcClientV1) CreateMap(ctx context.Context,
	map2d data1.Map2dV1) (*data1.Map2dV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"map", map2d,
	)

	response, err := c.CallCommand(ctx, "create_map", params)
	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.Map2dV1](response)
}

func (c *GeoRendererGrpcClientV1) UpdateMap(ctx context.Context,
	map2d data1.Map2dV1) (*data1.Map2dV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"map", map2d,
	)

	response, err := c.CallCommand(ctx, "update_map", params)
	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.Map2dV1](response)
}

func (c *GeoRendererGrpcClientV1) DeleteMapById(ctx context.Context,
	id string) (*data1.Map2dV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"map_id", id,
	)

	response, err := c.CallCommand(ctx, "delete_map_by_id", params)
	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.Map2dV1](response)
}
