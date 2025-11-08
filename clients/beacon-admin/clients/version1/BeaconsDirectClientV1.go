package clients1

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	logic "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/service"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	clients "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/clients"
)

type BeaconsDirectClientV1 struct {
	clients.DirectClient
	service logic.IBeaconsService
}

func NewBeaconsDirectClientV1() *BeaconsDirectClientV1 {
	c := &BeaconsDirectClientV1{
		DirectClient: *clients.NewDirectClient(),
	}
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("beacon-admin", "service", "*", "*", "1.0"))
	return c
}

func (c *BeaconsDirectClientV1) SetReferences(ctx context.Context, references cref.IReferences) {
	c.DirectClient.SetReferences(ctx, references)

	service, ok := c.Service.(logic.IBeaconsService)
	if !ok {
		panic("BeaconsDirectClientV1: Cant't resolv dependency 'service' to IBeaconsService")
	}
	c.service = service
}

func (c *BeaconsDirectClientV1) GetBeacons(ctx context.Context, reqctx cdata.RequestContextV1,
	filter *cquery.FilterParams, paging *cquery.PagingParams) (*cquery.DataPage[data1.BeaconV1], error) {
	timing := c.Instrument(ctx, "beacons.get_beacons")
	result, err := c.service.GetBeacons(ctx, reqctx, *filter, *paging)
	timing.EndTiming(ctx, err)
	return &result, err
}

func (c *BeaconsDirectClientV1) GetBeaconById(ctx context.Context, reqctx cdata.RequestContextV1,
	beaconId string) (*data1.BeaconV1, error) {
	timing := c.Instrument(ctx, "beacons.get_beacon_by_id")
	result, err := c.service.GetBeaconById(ctx, reqctx, beaconId)
	timing.EndTiming(ctx, err)
	return &result, err
}

func (c *BeaconsDirectClientV1) GetBeaconByUdi(ctx context.Context, reqctx cdata.RequestContextV1,
	udi string) (*data1.BeaconV1, error) {
	timing := c.Instrument(ctx, "beacons.get_beacon_by_udi")
	result, err := c.service.GetBeaconByUdi(ctx, reqctx, udi)
	timing.EndTiming(ctx, err)
	return &result, err
}

func (c *BeaconsDirectClientV1) CreateBeacon(ctx context.Context, reqctx cdata.RequestContextV1,
	beacon data1.BeaconV1) (*data1.BeaconV1, error) {
	timing := c.Instrument(ctx, "beacons.create_beacon")
	result, err := c.service.CreateBeacon(ctx, reqctx, beacon)
	timing.EndTiming(ctx, err)
	return &result, err
}

func (c *BeaconsDirectClientV1) UpdateBeacon(ctx context.Context, reqctx cdata.RequestContextV1,
	beacon data1.BeaconV1) (*data1.BeaconV1, error) {
	timing := c.Instrument(ctx, "beacons.update_beacon")
	result, err := c.service.UpdateBeacon(ctx, reqctx, beacon)
	timing.EndTiming(ctx, err)
	return &result, err
}

func (c *BeaconsDirectClientV1) DeleteBeaconById(ctx context.Context, reqctx cdata.RequestContextV1,
	beaconId string) (*data1.BeaconV1, error) {
	timing := c.Instrument(ctx, "beacons.delete_beacon_by_id")
	result, err := c.service.DeleteBeaconById(ctx, reqctx, beaconId)
	timing.EndTiming(ctx, err)
	return &result, err
}
