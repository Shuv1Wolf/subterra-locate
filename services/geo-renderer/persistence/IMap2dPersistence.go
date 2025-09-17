package persistence

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IMap2dPersistence interface {
	GetPageByFilter(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data1.Map2dV1], error)

	GetOneById(ctx context.Context, id string) (data1.Map2dV1, error)

	Create(ctx context.Context, item data1.Map2dV1) (data1.Map2dV1, error)

	Update(ctx context.Context, item data1.Map2dV1) (data1.Map2dV1, error)

	DeleteById(ctx context.Context, id string) (data1.Map2dV1, error)
}
