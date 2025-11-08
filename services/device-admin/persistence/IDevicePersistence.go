package persistence

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IDevicePersistence interface {
	GetPageByFilter(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.DeviceV1], error)

	GetOneById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.DeviceV1, error)

	Create(ctx context.Context, reqctx cdata.RequestContextV1, item data.DeviceV1) (data.DeviceV1, error)

	Update(ctx context.Context, reqctx cdata.RequestContextV1, item data.DeviceV1) (data.DeviceV1, error)

	DeleteById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (data.DeviceV1, error)
}
