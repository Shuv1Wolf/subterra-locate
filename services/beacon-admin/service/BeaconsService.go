package logic

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/persistence"
	"github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/publisher"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type BeaconsService struct {
	persistence   persistence.IBeaconsPersistence
	commandSet    *BeaconsCommandSet
	beaconsEvents publisher.IPublisher
}

func NewBeaconsService() *BeaconsService {
	c := &BeaconsService{}
	return c
}

func (c *BeaconsService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	// Read configuration parameters here...
}

func (c *BeaconsService) GetCommandSet() *ccmd.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewBeaconsCommandSet(c)
	}
	return &c.commandSet.CommandSet
}

func (c *BeaconsService) SetReferences(ctx context.Context, references cref.IReferences) {
	res, err := references.GetOneRequired(
		cref.NewDescriptor("beacon-admin", "persistence", "*", "*", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.persistence = res.(persistence.IBeaconsPersistence)

	res = references.GetOneOptional(
		cref.NewDescriptor("beacon-admin", "publisher", "nats", "beacons-events", "1.0"),
	)
	c.beaconsEvents = res.(publisher.IPublisher)
}

func (c *BeaconsService) GetBeacons(ctx context.Context, reqctx cdata.RequestContextV1,
	filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.BeaconV1], error) {
	return c.persistence.GetPageByFilter(ctx, reqctx, filter, paging)
}

func (c *BeaconsService) GetBeaconById(ctx context.Context, reqctx cdata.RequestContextV1,
	beaconId string) (data.BeaconV1, error) {

	return c.persistence.GetOneById(ctx, reqctx, beaconId)
}

func (c *BeaconsService) GetBeaconByUdi(ctx context.Context, reqctx cdata.RequestContextV1,
	beaconId string) (data.BeaconV1, error) {

	return c.persistence.GetOneByUdi(ctx, reqctx, beaconId)
}

func (c *BeaconsService) CreateBeacon(ctx context.Context, reqctx cdata.RequestContextV1,
	beacon data.BeaconV1) (data.BeaconV1, error) {

	if beacon.Type == "" {
		beacon.Type = data.Unknown
	}

	b, err := c.persistence.Create(ctx, reqctx, beacon)
	if err != nil {
		return b, err
	}

	if c.beaconsEvents != nil {
		err = c.beaconsEvents.SendBeaconCreatedEvent(ctx, b.Id)
		if err != nil {
			return b, err
		}
	}

	return b, nil
}

func (c *BeaconsService) UpdateBeacon(ctx context.Context, reqctx cdata.RequestContextV1,
	beacon data.BeaconV1) (data.BeaconV1, error) {

	if beacon.Type == "" {
		beacon.Type = data.Unknown
	}

	b, err := c.persistence.Update(ctx, reqctx, beacon)
	if err != nil {
		return b, err
	}

	if c.beaconsEvents != nil {
		err = c.beaconsEvents.SendBeaconChangedEvent(ctx, b.Id)
		if err != nil {
			return b, err
		}
	}

	return b, err
}

func (c *BeaconsService) DeleteBeaconById(ctx context.Context, reqctx cdata.RequestContextV1,
	beaconId string) (data.BeaconV1, error) {

	b, err := c.persistence.DeleteById(ctx, reqctx, beaconId)
	if err != nil {
		return b, err
	}

	if c.beaconsEvents != nil {
		err = c.beaconsEvents.SendBeaconCreatedEvent(ctx, b.Id)
		if err != nil {
			return b, err
		}
	}

	return b, err
}
