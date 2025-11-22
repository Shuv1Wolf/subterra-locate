package service

import (
	"context"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	natsConst "github.com/Shuv1Wolf/subterra-locate/services/common/nats/const"
	natsEvents "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/persistence"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/publisher"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type ZoneService struct {
	persistence persistence.IZonePersistence
	commandSet  *ZoneCommandSet
	zoneEvents  publisher.IPublisher
}

func NewZoneService() *ZoneService {
	c := &ZoneService{}
	return c
}

func (c *ZoneService) Configure(ctx context.Context, config *cconf.ConfigParams) {
}

func (c *ZoneService) GetCommandSet() *ccmd.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewZoneCommandSet(c)
	}
	return &c.commandSet.CommandSet
}

func (c *ZoneService) SetReferences(ctx context.Context, references cref.IReferences) {
	res, err := references.GetOneRequired(
		cref.NewDescriptor("zone-processor", "persistence", "*", "*", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.persistence = res.(persistence.IZonePersistence)

	res = references.GetOneOptional(
		cref.NewDescriptor("zone-processor", "publisher", "nats", "zone-events", "1.0"),
	)
	if res != nil {
		c.zoneEvents = res.(publisher.IPublisher)
	}
}

func (c *ZoneService) GetZones(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data1.ZoneV1], error) {
	return c.persistence.GetPageByFilter(ctx, reqctx, filter, paging)
}

func (c *ZoneService) GetZoneById(ctx context.Context, reqctx cdata.RequestContextV1, zone_id string) (data1.ZoneV1, error) {
	return c.persistence.GetOneById(ctx, reqctx, zone_id)
}

func (c *ZoneService) CreateZone(ctx context.Context, reqctx cdata.RequestContextV1, zone data1.ZoneV1) (data1.ZoneV1, error) {
	b, err := c.persistence.Create(ctx, reqctx, zone)
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

func (c *ZoneService) UpdateZone(ctx context.Context, reqctx cdata.RequestContextV1, zone data1.ZoneV1) (data1.ZoneV1, error) {
	b, err := c.persistence.Update(ctx, reqctx, zone)
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

	return b, nil
}

func (c *ZoneService) DeleteZoneById(ctx context.Context, reqctx cdata.RequestContextV1, zone_id string) (data1.ZoneV1, error) {
	b, err := c.persistence.DeleteById(ctx, reqctx, zone_id)
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

	return b, nil
}
