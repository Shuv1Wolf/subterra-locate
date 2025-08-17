package service

import (
	"context"
	"sync"

	"github.com/nats-io/nats.go"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	natsConst "github.com/Shuv1Wolf/subterra-locate/services/common/nats/const"

	bClient "github.com/Shuv1Wolf/subterra-locate/clients/beacon-admin/clients/version1"

	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/listener"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/publisher"
)

type LocationEngineService struct {
	Logger *clog.CompositeLogger

	rawBleListener          listener.IListener
	devicePositionPublisher publisher.IPublisher

	beaconsMap           map[string]*data1.BeaconV1
	beaconAdmin          bClient.IBeaconsClientV1
	beaconsEventListener listener.IListener

	mu     sync.RWMutex
	isOpen bool
}

func NewLocationEngineService() *LocationEngineService {
	return &LocationEngineService{
		Logger:     clog.NewCompositeLogger(),
		beaconsMap: map[string]*data1.BeaconV1{},
	}
}

func (c *LocationEngineService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.Logger.Configure(ctx, config)
}

func (c *LocationEngineService) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)

	res, err := references.GetOneRequired(
		cref.NewDescriptor("location-engine", "listener", "nats", "ble-raw-rssi", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.rawBleListener = res.(listener.IListener)

	res, err = references.GetOneRequired(
		cref.NewDescriptor("location-engine", "listener", "nats", "beacons-events", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.beaconsEventListener = res.(listener.IListener)

	res, err = references.GetOneRequired(
		cref.NewDescriptor("location-engine", "publisher", "nats", "device-position", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.devicePositionPublisher = res.(publisher.IPublisher)

	c.setClientsReferences(ctx, references)
}

func (c *LocationEngineService) setClientsReferences(ctx context.Context, references cref.IReferences) {
	res, err := references.GetOneRequired(
		cref.NewDescriptor("beacon-admin", "client", "*", "*", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.beaconAdmin = res.(bClient.IBeaconsClientV1)
}

func (c *LocationEngineService) Open(ctx context.Context) error {
	if c.isOpen {
		return nil
	}

	c.initBeaconsCache()

	c.Logger.Info(ctx, "Starting message listener for ble")
	if err := c.rawBleListener.Listen(ctx, c); err != nil {
		c.Logger.Error(ctx, err, "Error while listening to message bus")
	}

	c.Logger.Info(ctx, "Starting message listener for beacons")
	if err := c.beaconsEventListener.Listen(ctx, c); err != nil {
		c.Logger.Error(ctx, err, "Error while listening to message bus")
	}

	c.isOpen = true
	return nil
}

func (c *LocationEngineService) IsOpen() bool {
	return c.isOpen
}

func (c *LocationEngineService) Close(ctx context.Context) error {
	if c.isOpen {
		c.rawBleListener.EndListen(ctx)
		c.beaconsEventListener.EndListen(ctx)
		c.isOpen = false
	}
	return nil
}

func (c *LocationEngineService) ReceiveMessage(ctx context.Context, envelope *cqueues.MessageEnvelope, queue cqueues.IMessageQueue) error {
	var subject string

	if msg, ok := envelope.GetReference().(*nats.Msg); ok {
		subject = msg.Subject
	}

	switch subject {
	case natsConst.NATS_EVENTS_BLE_RSSI_TOPIC:
		c.bleEventHandler(ctx, envelope.GetMessageAsString())

	case natsConst.NATS_BEACONS_EVENTS_TOPIC:
		switch envelope.MessageType {
		case natsConst.NATS_BEACONS_EVENTS_CHANGED_TYPE, natsConst.NATS_BEACONS_EVENTS_CREATED_TYPE:
			c.beaconChangedEvent(ctx, envelope.GetMessageAsString())
		case natsConst.NATS_BEACONS_EVENTS_DELETED_TYPE:
			c.beaconDeletedEvent(ctx, envelope.GetMessageAsString())
		}
	default:
		c.Logger.Debug(ctx, "Unknown subject: "+subject)
	}
	return nil
}
