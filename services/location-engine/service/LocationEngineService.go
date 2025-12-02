package service

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	grpc "google.golang.org/grpc"

	bdata "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	natsConst "github.com/Shuv1Wolf/subterra-locate/services/common/nats/const"
	ddata "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"

	bclient "github.com/Shuv1Wolf/subterra-locate/clients/beacon-admin/clients/version1"
	dclient "github.com/Shuv1Wolf/subterra-locate/clients/device-admin/clients/version1"

	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/listener"
	protos "github.com/Shuv1Wolf/subterra-locate/services/location-engine/protos"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/publisher"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/utils"
)

type LocationEngineService struct {
	Logger *clog.CompositeLogger

	rawBleListener          listener.IListener
	devicePositionPublisher publisher.IPublisher

	beaconsMap           map[string]*bdata.BeaconV1
	beaconAdmin          bclient.IBeaconsClientV1
	beaconsEventListener listener.IListener

	deviceMap           map[string]*ddata.DeviceV1
	deviceAdmin         dclient.IDeviceClientV1
	deviceEventListener listener.IListener

	mu     sync.RWMutex
	isOpen bool

	monitorPort      string
	deviceStateStore *utils.DeviceStateStore
	beaconStateStore *utils.BeaconStateStore
}

func NewLocationEngineService() *LocationEngineService {
	return &LocationEngineService{
		Logger:           clog.NewCompositeLogger(),
		beaconsMap:       map[string]*bdata.BeaconV1{},
		deviceMap:        map[string]*ddata.DeviceV1{},
		deviceStateStore: utils.NewDeviceStateStore(),
		beaconStateStore: utils.NewBeaconStateStore(),
	}
}

func (c *LocationEngineService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.Logger.Configure(ctx, config)
	c.monitorPort = config.GetAsStringWithDefault("monitor.port", ":10030")
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
		cref.NewDescriptor("location-engine", "publisher", "nats", "device-position", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.devicePositionPublisher = res.(publisher.IPublisher)

	c.setEventsReferences(ctx, references)
	c.setClientsReferences(ctx, references)
}

func (c *LocationEngineService) setEventsReferences(ctx context.Context, references cref.IReferences) {
	res, err := references.GetOneRequired(
		cref.NewDescriptor("location-engine", "listener", "nats", "beacons-events", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.beaconsEventListener = res.(listener.IListener)

	res, err = references.GetOneRequired(
		cref.NewDescriptor("location-engine", "listener", "nats", "device-events", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.deviceEventListener = res.(listener.IListener)
}

func (c *LocationEngineService) setClientsReferences(ctx context.Context, references cref.IReferences) {
	res, err := references.GetOneRequired(
		cref.NewDescriptor("beacon-admin", "client", "*", "*", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.beaconAdmin = res.(bclient.IBeaconsClientV1)

	res, err = references.GetOneRequired(
		cref.NewDescriptor("device-admin", "client", "*", "*", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.deviceAdmin = res.(dclient.IDeviceClientV1)
}

func (c *LocationEngineService) Open(ctx context.Context) error {
	if c.isOpen {
		return nil
	}

	c.initBeaconsCache()
	c.initDeviceCache()
	c.startMessageListener(ctx)
	c.runMonitorLocation()
	c.startStaleDeviceChecker(ctx)

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
		c.deviceEventListener.EndListen(ctx)
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
		c.bleEventHandler(ctx, envelope)

	case natsConst.NATS_BEACONS_EVENTS_TOPIC:
		switch envelope.MessageType {
		case natsConst.NATS_BEACONS_EVENTS_CHANGED_TYPE, natsConst.NATS_BEACONS_EVENTS_CREATED_TYPE:
			c.beaconChangedEvent(ctx, envelope.GetMessageAsString())
		case natsConst.NATS_BEACONS_EVENTS_DELETED_TYPE:
			c.beaconDeletedEvent(ctx, envelope.GetMessageAsString())
		}

	case natsConst.NATS_DEVICE_EVENTS_TOPIC:
		switch envelope.MessageType {
		case natsConst.NATS_DEVICE_EVENTS_CHANGED_TYPE, natsConst.NATS_DEVICE_EVENTS_CREATED_TYPE:
			c.deviceChangedEvent(ctx, envelope.GetMessageAsString())
		case natsConst.NATS_DEVICE_EVENTS_DELETED_TYPE:
			c.deviceDeletedEvent(ctx, envelope.GetMessageAsString())
		}
	default:
		c.Logger.Debug(ctx, "Unknown subject: "+subject)
	}
	return nil
}

func (c *LocationEngineService) startMessageListener(ctx context.Context) {
	c.Logger.Info(ctx, "Starting message listener for ble")
	if err := c.rawBleListener.Listen(ctx, c); err != nil {
		c.Logger.Error(ctx, err, "Error while listening to message bus")
	}

	c.Logger.Info(ctx, "Starting message listener for beacons")
	if err := c.beaconsEventListener.Listen(ctx, c); err != nil {
		c.Logger.Error(ctx, err, "Error while listening to message bus")
	}

	c.Logger.Info(ctx, "Starting message listener for devices")
	if err := c.deviceEventListener.Listen(ctx, c); err != nil {
		c.Logger.Error(ctx, err, "Error while listening to message bus")
	}
}

func (c *LocationEngineService) startStaleDeviceChecker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute) // Check every minute
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.checkStaleDevices(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (c *LocationEngineService) checkStaleDevices(ctx context.Context) {
	devices := c.deviceStateStore.GetAllDevices()
	for _, device := range devices {
		if device.Online && time.Since(device.UpdatedAt) > 15*time.Minute {
			// TODO: добавить тревогу о том что девайс оффлайн
			c.Logger.Info(ctx, "Device %s is offline", device.DeviceID)
			device.X = 0
			device.Y = 0
			device.Z = 0
			device.Online = false
			device.UpdatedAt = time.Now()
			c.deviceStateStore.Upsert(device)
		}
	}
}

func (s *LocationEngineService) runMonitorLocation() {
	lis, err := net.Listen("tcp", s.monitorPort)
	if err != nil {
		s.Logger.Error(context.Background(), err, "Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	grpcServer := grpc.NewServer(opts...)

	monitorSvc := NewMonitorLocation(s.deviceStateStore, s.beaconStateStore, s.Logger)
	protos.RegisterLocationMonitorServer(grpcServer, monitorSvc)

	go func() {
		s.Logger.Info(context.Background(), "Starting monitor service on port %s", s.monitorPort)
		if err := grpcServer.Serve(lis); err != nil {
			s.Logger.Error(context.Background(), err, "Failed to serve: %v", err)
		}
	}()
}
