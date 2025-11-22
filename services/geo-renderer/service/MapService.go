package logic

import (
	"context"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	natsConst "github.com/Shuv1Wolf/subterra-locate/services/common/nats/const"
	natsEvents "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
	data "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/persistence"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/publisher"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type MapService struct {
	mapPersistence  persistence.IMap2dPersistence
	zonePersistence persistence.IZonePersistence

	commandSet *MapCommandSet

	map2dEvents publisher.IPublisher
	zoneEvents  publisher.IPublisher
}

func NewMapService() *MapService {
	c := &MapService{}
	return c
}

func (c *MapService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	// Read configuration parameters here...
}

func (c *MapService) GetCommandSet() *ccmd.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewMapCommandSet(c)
	}
	return &c.commandSet.CommandSet
}

func (c *MapService) SetReferences(ctx context.Context, references cref.IReferences) {
	res, err := references.GetOneRequired(
		cref.NewDescriptor("geo-renderer", "persistence", "*", "map-2d", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.mapPersistence = res.(persistence.IMap2dPersistence)

	res, err = references.GetOneRequired(
		cref.NewDescriptor("geo-renderer", "persistence", "*", "zone", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.zonePersistence = res.(persistence.IZonePersistence)

	res = references.GetOneOptional(
		cref.NewDescriptor("geo-renderer", "publisher", "nats", "map-2d-events", "1.0"),
	)
	if res != nil {
		c.map2dEvents = res.(publisher.IPublisher)
	}

	res = references.GetOneOptional(
		cref.NewDescriptor("geo-renderer", "publisher", "nats", "zone-events", "1.0"),
	)
	if res != nil {
		c.zoneEvents = res.(publisher.IPublisher)
	}
}

func (c *MapService) GetMaps(ctx context.Context, reqctx cdata.RequestContextV1,
	filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.Map2dV1], error) {
	return c.mapPersistence.GetPageByFilter(ctx, reqctx, filter, paging)
}

func (c *MapService) GetMapById(ctx context.Context, reqctx cdata.RequestContextV1,
	id string) (data.Map2dV1, error) {

	return c.mapPersistence.GetOneById(ctx, reqctx, id)
}

func (c *MapService) CreateMap(ctx context.Context, reqctx cdata.RequestContextV1,
	map2d data.Map2dV1) (data.Map2dV1, error) {
	b, err := c.mapPersistence.Create(ctx, reqctx, map2d)
	if err != nil {
		return b, err
	}

	if c.map2dEvents != nil {
		event := natsEvents.Map2dChangedEvent{
			Id: b.Id,
		}

		err = c.map2dEvents.SendEvent(ctx, event, natsConst.NATS_MAP2D_EVENTS_CREATED_TYPE)
		if err != nil {
			return b, err
		}
	}

	return b, nil
}

func (c *MapService) UpdateMap(ctx context.Context, reqctx cdata.RequestContextV1,
	map2d data.Map2dV1) (data.Map2dV1, error) {

	b, err := c.mapPersistence.Update(ctx, reqctx, map2d)
	if err != nil {
		return b, err
	}

	if c.map2dEvents != nil {
		event := natsEvents.Map2dChangedEvent{
			Id: b.Id,
		}

		err = c.map2dEvents.SendEvent(ctx, event, natsConst.NATS_MAP2D_EVENTS_CHANGED_TYPE)
		if err != nil {
			return b, err
		}
	}

	return b, err
}

func (c *MapService) DeleteMapById(ctx context.Context, reqctx cdata.RequestContextV1,
	id string) (data.Map2dV1, error) {

	b, err := c.mapPersistence.DeleteById(ctx, reqctx, id)
	if err != nil {
		return b, err
	}

	if c.map2dEvents != nil {
		event := natsEvents.Map2dChangedEvent{
			Id: b.Id,
		}

		err = c.map2dEvents.SendEvent(ctx, event, natsConst.NATS_MAP2D_EVENTS_DELETED_TYPE)
		if err != nil {
			return b, err
		}
	}

	return b, err
}

func (c *MapService) GetZones(ctx context.Context, reqctx cdata.RequestContextV1,
	filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.ZoneV1], error) {
	return c.zonePersistence.GetPageByFilter(ctx, reqctx, filter, paging)
}

func (c *MapService) GetZoneById(ctx context.Context, reqctx cdata.RequestContextV1,
	id string) (data.ZoneV1, error) {

	return c.zonePersistence.GetOneById(ctx, reqctx, id)
}

func (c *MapService) CreateZone(ctx context.Context, reqctx cdata.RequestContextV1,
	zone data.ZoneV1) (data.ZoneV1, error) {
	b, err := c.zonePersistence.Create(ctx, reqctx, zone)
	if err != nil {
		return b, err
	}

	if c.zoneEvents != nil {
		event := natsEvents.ZoneChangedEvent{
			Id: b.Id,
		}

		err = c.zoneEvents.SendEvent(ctx, event, natsConst.NATS_ZONE_EVENTS_CREATED_TYPE)
		if err != nil {
			return b, err
		}
	}

	return b, nil
}

func (c *MapService) UpdateZone(ctx context.Context, reqctx cdata.RequestContextV1,
	zone data.ZoneV1) (data.ZoneV1, error) {

	b, err := c.zonePersistence.Update(ctx, reqctx, zone)
	if err != nil {
		return b, err
	}

	if c.zoneEvents != nil {
		event := natsEvents.ZoneChangedEvent{
			Id: b.Id,
		}

		err = c.zoneEvents.SendEvent(ctx, event, natsConst.NATS_ZONE_EVENTS_CHANGED_TYPE)
		if err != nil {
			return b, err
		}
	}

	return b, err
}

func (c *MapService) DeleteZoneById(ctx context.Context, reqctx cdata.RequestContextV1,
	id string) (data.ZoneV1, error) {

	b, err := c.zonePersistence.DeleteById(ctx, reqctx, id)
	if err != nil {
		return b, err
	}

	if c.zoneEvents != nil {
		event := natsEvents.ZoneChangedEvent{
			Id: b.Id,
		}

		err = c.zoneEvents.SendEvent(ctx, event, natsConst.NATS_ZONE_EVENTS_DELETED_TYPE)
		if err != nil {
			return b, err
		}
	}

	return b, err
}
