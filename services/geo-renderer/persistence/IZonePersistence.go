package persistence

import (
	"context"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IZonePersistence interface {
	GetPageByFilter(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data1.ZoneV1], error)

	GetOneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data1.ZoneV1, error)

	Create(ctx context.Context, reqctx cdata.RequestContextV1, item data1.ZoneV1) (data1.ZoneV1, error)

	Update(ctx context.Context, reqctx cdata.RequestContextV1, item data1.ZoneV1) (data1.ZoneV1, error)

	DeleteById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data1.ZoneV1, error)
}
