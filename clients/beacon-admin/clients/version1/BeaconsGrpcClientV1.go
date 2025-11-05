package clients1

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
	clients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
)

type BeaconsGrpcClientV1 struct {
	*cclients.CommandableGrpcClient
}

func NewBeaconsGrpcClientV1() *BeaconsGrpcClientV1 {
	c := &BeaconsGrpcClientV1{
		CommandableGrpcClient: cclients.NewCommandableGrpcClient("beacon.admin.v1"),
	}
	return c
}

func (c *BeaconsGrpcClientV1) GetBeacons(ctx context.Context,
	filter *cquery.FilterParams,
	paging *cquery.PagingParams) (*cquery.DataPage[data1.BeaconV1], error) {

	var pagingMap map[string]interface{}

	if paging != nil {
		pagingMap = map[string]interface{}{
			"skip":  paging.Skip,
			"take":  paging.Take,
			"total": paging.Total,
		}
	}

	params := cdata.NewAnyValueMapFromTuples(
		"filter", filter.StringValueMap.Value(),
		"paging", pagingMap,
	)

	response, err := c.CallCommand(ctx, "get_beacons", cdata.NewAnyValueMapFromValue(params.Value()))

	if err != nil {
		return cquery.NewEmptyDataPage[data1.BeaconV1](), err
	}

	return clients.HandleHttpResponse[*cquery.DataPage[data1.BeaconV1]](response)
}

func (c *BeaconsGrpcClientV1) GetBeaconById(ctx context.Context,
	beaconId string) (*data1.BeaconV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"beacon_id", beaconId,
	)

	response, err := c.CallCommand(ctx, "get_beacon_by_id", params)

	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[*data1.BeaconV1](response)
}

func (c *BeaconsGrpcClientV1) GetBeaconByUdi(ctx context.Context,
	udi string) (*data1.BeaconV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"udi", udi,
	)

	response, err := c.CallCommand(ctx, "get_beacon_by_udi", params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[*data1.BeaconV1](response)
}

func (c *BeaconsGrpcClientV1) CreateBeacon(ctx context.Context,
	beacon data1.BeaconV1) (*data1.BeaconV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"beacon", beacon,
	)

	response, err := c.CallCommand(ctx, "create_beacon", params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[*data1.BeaconV1](response)
}

func (c *BeaconsGrpcClientV1) UpdateBeacon(ctx context.Context,
	beacon data1.BeaconV1) (*data1.BeaconV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"beacon", beacon,
	)

	response, err := c.CallCommand(ctx, "update_beacon", params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[*data1.BeaconV1](response)
}

func (c *BeaconsGrpcClientV1) DeleteBeaconById(ctx context.Context,
	beaconId string) (*data1.BeaconV1, error) {
	params := cdata.NewAnyValueMapFromTuples(
		"beacon_id", beaconId,
	)

	response, err := c.CallCommand(ctx, "delete_beacon_by_id", params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[*data1.BeaconV1](response)
}
