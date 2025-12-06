package service

import (
	"context"
	"encoding/json"
	"net"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	natsConst "github.com/Shuv1Wolf/subterra-locate/services/common/nats/const"
	natsEvents "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/listener"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/persistence"
	protos "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/protos"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/publisher"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/utils"
	"github.com/nats-io/nats.go"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
	"google.golang.org/grpc"
)

type ZoneService struct {
	persistence persistence.IZonePersistence
	commandSet  *ZoneCommandSet

	zoneEvents publisher.IPublisher
	listener   listener.IListener

	logger *clog.CompositeLogger

	stateStore  *utils.ZoneStateStore
	processor   *ZoneStateProcessor
	monitorPort string
	isOpen      bool
}

func NewZoneService() *ZoneService {
	store := utils.NewZoneStateStore()
	c := &ZoneService{
		logger:     clog.NewCompositeLogger(),
		stateStore: store,
		processor:  NewZoneStateProcessor(store),
	}
	return c
}

func (c *ZoneService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.logger.Configure(ctx, config)
	c.monitorPort = config.GetAsStringWithDefault("monitor.port", ":10070")
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
	c.logger.SetReferences(ctx, references)

	res, err = references.GetOneRequired(
		cref.NewDescriptor("zone-processor", "listener", "nats", "device-position", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.listener = res.(listener.IListener)
}

func (c *ZoneService) Open(ctx context.Context) error {
	c.initCache(ctx)
	c.runMonitorLocation()

	c.logger.Info(ctx, "Starting message listener")
	if err := c.listener.Listen(ctx, c); err != nil {
		c.logger.Error(ctx, err, "Error while listening to message bus")
	}

	c.isOpen = true
	return nil
}

func (c *ZoneService) Close(ctx context.Context) error {
	c.isOpen = false
	return nil
}

func (c *ZoneService) IsOpen() bool {
	return c.isOpen
}

func (c *ZoneService) initCache(ctx context.Context) {
	limit := int64(100)
	skip := int64(0)

	for {
		page := *cquery.NewPagingParams(skip, limit, false)
		zones, err := c.persistence.GetPageByFilter(ctx, cdata.RequestContextV1{}, *cquery.NewEmptyFilterParams(), page)
		if err != nil {
			c.logger.Error(ctx, err, "Error getting zones")
			return
		}

		if len(zones.Data) == 0 {
			break
		}

		for _, zone := range zones.Data {
			c.stateStore.Upsert(&zone)
		}

		if int64(len(zones.Data)) < limit {
			break
		}

		skip += limit
	}

	c.logger.Info(context.Background(), "Zones stored in cache")
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

	c.stateStore.Upsert(&b)
	c.processor.HandleZoneAdd(&b)

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

	c.stateStore.Upsert(&b)
	c.processor.HandleZoneUpdate(&b)

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
	c.stateStore.Delete(&b)
	c.processor.HandleZoneDelete(&b)

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

func (c *ZoneService) ReceiveMessage(ctx context.Context, envelope *cqueues.MessageEnvelope, queue cqueues.IMessageQueue) error {
	var subject string

	if msg, ok := envelope.GetReference().(*nats.Msg); ok {
		subject = msg.Subject
	}

	switch subject {
	case natsConst.NATS_EVENTS_DEVICE_POSITION:
		c.hendleDevicePosition(ctx, envelope.GetMessageAsString())
	default:
		c.logger.Debug(ctx, "Unknown subject: "+subject)
	}
	return nil
}

func (c *ZoneService) hendleDevicePosition(ctx context.Context, msg string) {
	var event natsEvents.DevicePositioningEventV1
	err := json.Unmarshal([]byte(msg), &event)
	if err != nil {
		c.logger.Error(ctx, err, "Failed to deserialize message")
		return
	}

	c.processor.ProcessPosition(event)
}

func (s *ZoneService) runMonitorLocation() {
	lis, err := net.Listen("tcp", s.monitorPort)
	if err != nil {
		s.logger.Error(context.Background(), err, "Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	grpcServer := grpc.NewServer(opts...)

	monitorSvc := NewZoneMonitorService(s.stateStore, s.logger)
	protos.RegisterZoneMonitorServer(grpcServer, monitorSvc)

	go func() {
		s.logger.Info(context.Background(), "Starting monitor service on port %s", s.monitorPort)
		if err := grpcServer.Serve(lis); err != nil {
			s.logger.Error(context.Background(), err, "Failed to serve: %v", err)
		}
	}()
}
