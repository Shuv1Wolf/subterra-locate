package service

import (
	"context"
	"encoding/json"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"

	mqtt "github.com/Shuv1Wolf/subterra-locate/services/common/mqtt"
	natsEvents "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"

	"github.com/Shuv1Wolf/subterra-locate/services/mqtt-listener/listener"
	"github.com/Shuv1Wolf/subterra-locate/services/mqtt-listener/publisher"
)

type MqttListenerService struct {
	Logger          *clog.CompositeLogger
	bleRssiListener listener.IMqttListener
	natsPublisher   publisher.IPublisher
	isOpen          bool
}

func NewMqttListenerService() *MqttListenerService {
	return &MqttListenerService{
		Logger: clog.NewCompositeLogger(),
	}
}

func (c *MqttListenerService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.Logger.Configure(ctx, config)
}

func (c *MqttListenerService) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)

	res, err := references.GetOneRequired(
		cref.NewDescriptor("mqtt-listener", "listener", "mqtt", "ble-rssi", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.bleRssiListener = res.(listener.IMqttListener)

	res, err = references.GetOneRequired(
		cref.NewDescriptor("mqtt-listener", "publisher", "nats", "loc-raw-ble", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.natsPublisher = res.(publisher.IPublisher)
}

func (c *MqttListenerService) Open(ctx context.Context) error {
	if c.isOpen {
		return nil
	}
	c.isOpen = true
	c.startMessageListener(ctx)
	return nil
}

func (c *MqttListenerService) IsOpen() bool {
	return c.isOpen
}

func (c *MqttListenerService) Close(ctx context.Context) error {
	if c.isOpen {
		c.bleRssiListener.EndListen(ctx)
		c.isOpen = false
	}
	return nil
}

func (c *MqttListenerService) startMessageListener(ctx context.Context) {
	go func() {
		c.Logger.Info(ctx, "Starting message listener...")

		if err := c.bleRssiListener.Listen(ctx, c); err != nil {
			c.Logger.Error(ctx, err, "Error while listening to message bus")
		}
	}()
}

func (c *MqttListenerService) ReceiveMessage(ctx context.Context, envelope *cqueues.MessageEnvelope, queue cqueues.IMessageQueue) error {
	switch envelope.MessageType {
	case mqtt.MQTT_BLE_RSSI_TOPIC:
		var event natsEvents.BLEBeaconRawEventV1
		err := json.Unmarshal([]byte(envelope.GetMessageAsString()), &event)
		if err != nil {
			c.Logger.Error(ctx, err, "Failed to deserialize message")
		}

		err = c.natsPublisher.SendRawBle(ctx, &event)
		if err != nil {
			c.Logger.Error(ctx, err, "Failed to send message")
		}
		c.Logger.Debug(ctx, "Message sent: "+envelope.MessageType)

	default:
		c.Logger.Debug(ctx, "Unknown message type: "+envelope.MessageType)
	}
	return nil
}
