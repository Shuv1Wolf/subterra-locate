package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	natsConst "github.com/Shuv1Wolf/subterra-locate/services/common/nats/const"
	natsEvents "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
	"github.com/nats-io/nats.go"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	keys "github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-history/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-history/listener"
	pers "github.com/Shuv1Wolf/subterra-locate/services/geo-history/persistence"
)

const ( // TODO: Move to config
	batchSize   = 100
	flushPeriod = 5 * time.Second
)

type GeoHistoryService struct {
	Logger *clog.CompositeLogger

	persistence pers.IGeoHistoryPersistence
	commandSet  *GeoHistoryCommandSet
	geoListener listener.IListener

	isOpen bool

	mu       sync.Mutex
	batchBuf []*data1.HistoricalRecordV1
	stopChan chan struct{}
}

func NewGeoHistoryService() *GeoHistoryService {
	c := &GeoHistoryService{
		Logger: clog.NewCompositeLogger(),
	}
	return c
}

func (c *GeoHistoryService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.Logger.Configure(ctx, config)
}

func (c *GeoHistoryService) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)
	res, err := references.GetOneRequired(
		cref.NewDescriptor("geo-history", "persistence", "*", "device", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.persistence = res.(pers.IGeoHistoryPersistence)

	res, err = references.GetOneRequired(
		cref.NewDescriptor("geo-history", "listener", "nats", "device-position", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.geoListener = res.(listener.IListener)
}

func (c *GeoHistoryService) IsOpen() bool {
	return c.isOpen
}

func (c *GeoHistoryService) Open(ctx context.Context) error {
	if c.isOpen {
		return nil
	}

	c.stopChan = make(chan struct{})
	c.batchBuf = make([]*data1.HistoricalRecordV1, 0, batchSize)

	go c.autoFlush(ctx)

	c.Logger.Info(ctx, "Starting message listener")
	if err := c.geoListener.Listen(ctx, c); err != nil {
		c.Logger.Error(ctx, err, "Error while listening to message bus")
	}

	c.isOpen = true
	return nil
}

func (c *GeoHistoryService) Close(ctx context.Context) error {
	if c.isOpen {
		close(c.stopChan)
		c.flush(ctx)
		c.geoListener.EndListen(ctx)
		c.isOpen = false
	}
	return nil
}

func (c *GeoHistoryService) GetCommandSet() *ccmd.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewGeoHistoryCommandSet(c)
	}
	return &c.commandSet.CommandSet
}

func (c *GeoHistoryService) GetHistory(ctx context.Context, reqctx cdata.RequestContextV1, mapId, from, to string, paging cquery.PagingParams, sortField cquery.SortField) (cquery.DataPage[data1.HistoricalRecordV1], error) {
	filter := cquery.NewFilterParamsFromTuples(
		"map_id", mapId,
		"from", from,
		"to", to,
	)

	return c.persistence.GetHistory(ctx, reqctx, *filter, paging, sortField)
}

func (c *GeoHistoryService) ReceiveMessage(ctx context.Context, envelope *cqueues.MessageEnvelope, queue cqueues.IMessageQueue) error {
	var subject string

	if msg, ok := envelope.GetReference().(*nats.Msg); ok {
		subject = msg.Subject
	}

	switch subject {
	case natsConst.NATS_EVENTS_DEVICE_POSITION:
		c.hendleDevicePosition(ctx, envelope.GetMessageAsString())
	default:
		c.Logger.Debug(ctx, "Unknown subject: "+subject)
	}
	return nil
}

func (c *GeoHistoryService) hendleDevicePosition(ctx context.Context, msg string) {
	var event natsEvents.DevicePositioningEventV1
	err := json.Unmarshal([]byte(msg), &event)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to deserialize message")
	}

	record := &data1.HistoricalRecordV1{
		Id:        keys.IdGenerator.NextLong(),
		OrgId:     event.OrgId,
		MapId:     event.MapId,
		EntityId:  event.DeviceId,
		Timestamp: event.Timestamp,
		X:         event.X,
		Y:         event.Y,
		Z:         event.Z,
	}

	c.mu.Lock()
	c.batchBuf = append(c.batchBuf, record)
	shouldFlush := len(c.batchBuf) >= batchSize
	c.mu.Unlock()

	if shouldFlush {
		c.flush(ctx)
	}
}

func (c *GeoHistoryService) autoFlush(ctx context.Context) {
	ticker := time.NewTicker(flushPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.flush(ctx)
		case <-c.stopChan:
			return
		}
	}
}

func (c *GeoHistoryService) flush(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.batchBuf) == 0 {
		return
	}

	toInsert := make([]*data1.HistoricalRecordV1, len(c.batchBuf))
	copy(toInsert, c.batchBuf)
	c.batchBuf = c.batchBuf[:0]

	if err := c.persistence.InsertBatch(ctx, toInsert); err != nil {
		c.Logger.Error(ctx, err, "Failed to insert batch of historical records")
	} else {
		c.Logger.Info(ctx, "Inserted batch of %d records", len(toInsert))
	}
}
