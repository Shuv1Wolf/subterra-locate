package clients1

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
)

type DeviceGrpcClientV1 struct {
	*cclients.CommandableGrpcClient
}

func NewDeviceGrpcClientV1() *DeviceGrpcClientV1 {
	c := &DeviceGrpcClientV1{
		CommandableGrpcClient: cclients.NewCommandableGrpcClient("device.admin.v1"),
	}
	return c
}

func (c *DeviceGrpcClientV1) GetDevice(ctx context.Context,
	filter cquery.FilterParams,
	paging cquery.PagingParams) (*cquery.DataPage[data1.DeviceV1], error) {

	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, &filter)
	c.AddPagingParams(params, &paging)

	response, err := c.CallCommand(ctx, "get_device", cdata.NewAnyValueMapFromValue(params.Value()))

	if err != nil {
		return cquery.NewEmptyDataPage[data1.DeviceV1](), err
	}

	return cclients.HandleHttpResponse[*cquery.DataPage[data1.DeviceV1]](response)
}

func (c *DeviceGrpcClientV1) GetDeviceById(ctx context.Context,
	id string) (*data1.DeviceV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"device_id", id,
	)

	response, err := c.CallCommand(ctx, "get_device_by_id", params)

	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.DeviceV1](response)
}

func (c *DeviceGrpcClientV1) CreateDevice(ctx context.Context,
	device data1.DeviceV1) (*data1.DeviceV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"device", device,
	)

	response, err := c.CallCommand(ctx, "create_device", params)
	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.DeviceV1](response)
}

func (c *DeviceGrpcClientV1) UpdateDevice(ctx context.Context,
	device data1.DeviceV1) (*data1.DeviceV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"device", device,
	)

	response, err := c.CallCommand(ctx, "update_device", params)
	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.DeviceV1](response)
}

func (c *DeviceGrpcClientV1) DeleteDeviceById(ctx context.Context,
	id string) (*data1.DeviceV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"device_id", id,
	)

	response, err := c.CallCommand(ctx, "delete_device_by_id", params)
	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.DeviceV1](response)
}
