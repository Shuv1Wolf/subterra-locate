package persistence

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IBeaconsPersistence interface {
	GetPageByFilter(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.BeaconV1], error)

	GetOneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.BeaconV1, error)

	GetOneByUdi(ctx context.Context, reqctx cdata.RequestContextV1, udi string) (data.BeaconV1, error)

	Create(ctx context.Context, reqctx cdata.RequestContextV1, item data.BeaconV1) (data.BeaconV1, error)

	Update(ctx context.Context, reqctx cdata.RequestContextV1, item data.BeaconV1) (data.BeaconV1, error)

	DeleteById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.BeaconV1, error)
}
