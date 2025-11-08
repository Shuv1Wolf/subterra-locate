package clients1

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IDeviceClientV1 interface {
	GetDevices(ctx context.Context, reqctx cdata.RequestContextV1, filter *cquery.FilterParams,
		paging *cquery.PagingParams) (*cquery.DataPage[data1.DeviceV1], error)

	GetDeviceById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (*data1.DeviceV1, error)

	CreateDevice(ctx context.Context, reqctx cdata.RequestContextV1, device data1.DeviceV1) (*data1.DeviceV1, error)

	UpdateDevice(ctx context.Context, reqctx cdata.RequestContextV1, device data1.DeviceV1) (*data1.DeviceV1, error)

	DeleteDeviceById(ctx context.Context, reqctx cdata.RequestContextV1, id string) (*data1.DeviceV1, error)
}
