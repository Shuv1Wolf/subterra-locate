package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"

	natsConst "github.com/Shuv1Wolf/subterra-locate/services/common/nats/const"
	natsEvents "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/listener"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/publisher"
)

type LocationEngineService struct {
	Logger         *clog.CompositeLogger
	rawBleListener listener.IListener
	publisher      publisher.IPublisher
	isOpen         bool
	// TODO: beacon cache
}

func NewLocationEngineService() *LocationEngineService {
	return &LocationEngineService{
		Logger: clog.NewCompositeLogger(),
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
		cref.NewDescriptor("location-engine", "publisher", "nats", "device-position", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.publisher = res.(publisher.IPublisher)
}

func (c *LocationEngineService) Open(ctx context.Context) error {
	if c.isOpen {
		return nil
	}
	c.isOpen = true
	c.startMessageListener(ctx)
	return nil
}

func (c *LocationEngineService) IsOpen() bool {
	return c.isOpen
}

func (c *LocationEngineService) Close(ctx context.Context) error {
	if c.isOpen {
		c.rawBleListener.EndListen(ctx)
		c.isOpen = false
	}
	return nil
}

func (c *LocationEngineService) startMessageListener(ctx context.Context) {
	go func() {
		c.Logger.Info(ctx, "Starting message listener...")

		if err := c.rawBleListener.Listen(ctx, c); err != nil {
			c.Logger.Error(ctx, err, "Error while listening to message bus")
		}
	}()
}

func (c *LocationEngineService) ReceiveMessage(ctx context.Context, envelope *cqueues.MessageEnvelope, queue cqueues.IMessageQueue) error {
	var subject string

	if msg, ok := envelope.GetReference().(*nats.Msg); ok {
		subject = msg.Subject
	}

	switch subject {
	case natsConst.NATS_EVENTS_BLE_RSSI_TOPIC:
		var event natsEvents.DeviceDetectedBLERawEventV1
		err := json.Unmarshal([]byte(envelope.GetMessageAsString()), &event)
		if err != nil {
			c.Logger.Error(ctx, err, "Failed to deserialize message")
		}
		fmt.Println(event)
		// TODO: process

		err = c.publisher.SendDevicePosition(ctx, &natsEvents.DevicePositioningEventV1{Test: "test"})
		if err != nil {
			c.Logger.Error(ctx, err, "Failed to send message")
		}

	default:
		c.Logger.Debug(ctx, "Unknown subject: "+subject)
	}
	return nil
}
