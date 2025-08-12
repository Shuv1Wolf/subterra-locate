package build

import (
	"github.com/Shuv1Wolf/subterra-locate/services/mqtt-listener/listener"
	"github.com/Shuv1Wolf/subterra-locate/services/mqtt-listener/publisher"
	"github.com/Shuv1Wolf/subterra-locate/services/mqtt-listener/service"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type MqttListenerServiceFactory struct {
	*cbuild.Factory
}

func NewMqttListenerServiceFactory() *MqttListenerServiceFactory {
	c := MqttListenerServiceFactory{}
	c.Factory = cbuild.NewFactory()

	mqttBleRssiListenerDescriptor := cref.NewDescriptor("mqtt-listener", "listener", "mqtt", "ble-rssi", "1.0")
	natsBleRawPublisherDescriptor := cref.NewDescriptor("mqtt-listener", "publisher", "nats", "loc-raw-ble", "1.0")
	serviceDescriptor := cref.NewDescriptor("mqtt-listener", "service", "default", "*", "1.0")

	c.RegisterType(mqttBleRssiListenerDescriptor, listener.NewMqttListener)
	c.RegisterType(serviceDescriptor, service.NewMqttListenerService)
	c.RegisterType(natsBleRawPublisherDescriptor, publisher.NewNatsPublisher)

	return &c
}
