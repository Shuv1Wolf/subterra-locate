package clients1

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type DeviceNullClientV1 struct {
}

func NewDeviceNullClientV1() *DeviceNullClientV1 {
	return &DeviceNullClientV1{}
}

func (c *DeviceNullClientV1) GetDevices(ctx context.Context, reqctx cdata.RequestContextV1,
	filter *cquery.FilterParams,
	paging *cquery.PagingParams) (*cquery.DataPage[data1.DeviceV1], error) {
	return cquery.NewEmptyDataPage[data1.DeviceV1](), nil
}

func (c *DeviceNullClientV1) GetDeviceById(ctx context.Context, reqctx cdata.RequestContextV1,
	id string) (*data1.DeviceV1, error) {
	return nil, nil
}

func (c *DeviceNullClientV1) CreateDevice(ctx context.Context, reqctx cdata.RequestContextV1,
	device data1.DeviceV1) (*data1.DeviceV1, error) {
	return nil, nil
}

func (c *DeviceNullClientV1) UpdateDevice(ctx context.Context, reqctx cdata.RequestContextV1,
	device data1.DeviceV1) (*data1.DeviceV1, error) {
	return nil, nil
}

func (c *DeviceNullClientV1) DeleteDeviceById(ctx context.Context, reqctx cdata.RequestContextV1,
	id string) (*data1.DeviceV1, error) {
	return nil, nil
}
