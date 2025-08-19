package logic

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IDeviceService interface {
	GetDevices(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.DeviceV1], error)

	GetDeviceById(ctx context.Context, deviceId string) (data.DeviceV1, error)

	CreateDevice(ctx context.Context, device data.DeviceV1) (data.DeviceV1, error)

	UpdateDevice(ctx context.Context, device data.DeviceV1) (data.DeviceV1, error)

	DeleteDeviceById(ctx context.Context, deviceId string) (data.DeviceV1, error)
}
