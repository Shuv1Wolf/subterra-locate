package service

import (
	"context"
	"fmt"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"

	mqttconst "github.com/Shuv1Wolf/subterra-locate/services/common/mqtt-const"

	"github.com/Shuv1Wolf/subterra-locate/services/mqtt-listener/listener"
)

type MqttListenerService struct {
	Logger          *clog.CompositeLogger
	bleRssiListener listener.IMqttListener
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
		cref.NewDescriptor("mqtt-listener", "listener", "*", "ble-rssi", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.bleRssiListener = res.(listener.IMqttListener)

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
	case mqttconst.MQTT_BLE_RSSI_TOPIC:
		// TODO: Handle beacon
		fmt.Println(envelope.GetMessageAsString())
	default:
		c.Logger.Debug(ctx, "Unknown message type: "+envelope.MessageType)
	}
	return nil
}
