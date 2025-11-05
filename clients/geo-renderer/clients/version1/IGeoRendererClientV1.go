package clients1

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IGeoRendererClientV1 interface {
	GetMaps(ctx context.Context, filter *cquery.FilterParams, paging *cquery.PagingParams) (*cquery.DataPage[data1.Map2dV1], error)

	GetMapById(ctx context.Context, id string) (*data1.Map2dV1, error)

	CreateMap(ctx context.Context, map2d data1.Map2dV1) (*data1.Map2dV1, error)

	UpdateMap(ctx context.Context, map2d data1.Map2dV1) (*data1.Map2dV1, error)

	DeleteMapById(ctx context.Context, id string) (*data1.Map2dV1, error)
}
