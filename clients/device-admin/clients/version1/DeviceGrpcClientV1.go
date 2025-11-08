package clients1

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	rdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
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

func (c *DeviceGrpcClientV1) GetDevices(ctx context.Context, reqctx rdata.RequestContextV1,
	filter *cquery.FilterParams,
	paging *cquery.PagingParams) (*cquery.DataPage[data1.DeviceV1], error) {

	var pagingMap map[string]interface{}
	var reqctxMap map[string]interface{}

	if paging != nil {
		pagingMap = map[string]interface{}{
			"skip":  paging.Skip,
			"take":  paging.Take,
			"total": paging.Total,
		}
	}

	reqctxMap = map[string]interface{}{
		"org_id":  reqctx.OrgId,
		"user_id": reqctx.UserId,
	}

	params := cdata.NewAnyValueMapFromTuples(
		"filter", filter.StringValueMap.Value(),
		"paging", pagingMap,
		"reqctx", reqctxMap,
	)

	response, err := c.CallCommand(ctx, "get_devices", cdata.NewAnyValueMapFromValue(params.Value()))

	if err != nil {
		return cquery.NewEmptyDataPage[data1.DeviceV1](), err
	}

	return cclients.HandleHttpResponse[*cquery.DataPage[data1.DeviceV1]](response)
}

func (c *DeviceGrpcClientV1) GetDeviceById(ctx context.Context, reqctx rdata.RequestContextV1,
	id string) (*data1.DeviceV1, error) {

	var reqctxMap map[string]interface{}
	reqctxMap = map[string]interface{}{
		"org_id":  reqctx.OrgId,
		"user_id": reqctx.UserId,
	}

	params := cdata.NewAnyValueMapFromTuples(
		"device_id", id,
		"reqctx", reqctxMap,
	)

	response, err := c.CallCommand(ctx, "get_device_by_id", params)

	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.DeviceV1](response)
}

func (c *DeviceGrpcClientV1) CreateDevice(ctx context.Context, reqctx rdata.RequestContextV1,
	device data1.DeviceV1) (*data1.DeviceV1, error) {

	var reqctxMap map[string]interface{}
	reqctxMap = map[string]interface{}{
		"org_id":  reqctx.OrgId,
		"user_id": reqctx.UserId,
	}

	params := cdata.NewAnyValueMapFromTuples(
		"device", device,
		"reqctx", reqctxMap,
	)

	response, err := c.CallCommand(ctx, "create_device", params)
	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.DeviceV1](response)
}

func (c *DeviceGrpcClientV1) UpdateDevice(ctx context.Context, reqctx rdata.RequestContextV1,
	device data1.DeviceV1) (*data1.DeviceV1, error) {

	var reqctxMap map[string]interface{}
	reqctxMap = map[string]interface{}{
		"org_id":  reqctx.OrgId,
		"user_id": reqctx.UserId,
	}

	params := cdata.NewAnyValueMapFromTuples(
		"device", device,
		"reqctx", reqctxMap,
	)

	response, err := c.CallCommand(ctx, "update_device", params)
	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.DeviceV1](response)
}

func (c *DeviceGrpcClientV1) DeleteDeviceById(ctx context.Context, reqctx rdata.RequestContextV1,
	id string) (*data1.DeviceV1, error) {

	var reqctxMap map[string]interface{}
	reqctxMap = map[string]interface{}{
		"org_id":  reqctx.OrgId,
		"user_id": reqctx.UserId,
	}

	params := cdata.NewAnyValueMapFromTuples(
		"device_id", id,
		"reqctx", reqctxMap,
	)

	response, err := c.CallCommand(ctx, "delete_device_by_id", params)
	if err != nil {
		return nil, err
	}

	return cclients.HandleHttpResponse[*data1.DeviceV1](response)
}
