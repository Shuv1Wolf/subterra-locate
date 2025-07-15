package persistence

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IBeaconsPersistence interface {
	GetPageByFilter(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.BeaconV1], error)

	GetOneById(ctx context.Context, id string) (data.BeaconV1, error)

	GetOneByUdi(ctx context.Context, udi string) (data.BeaconV1, error)

	Create(ctx context.Context, item data.BeaconV1) (data.BeaconV1, error)

	Update(ctx context.Context, item data.BeaconV1) (data.BeaconV1, error)

	DeleteById(ctx context.Context, id string) (data.BeaconV1, error)
}
