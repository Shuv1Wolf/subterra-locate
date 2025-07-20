package build

import (
	"github.com/Shuv1Wolf/subterra-locate/services/mqtt-listener/messagebus"
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

	mqttMessageBusDescriptor := cref.NewDescriptor("mqtt-listener", "messagebus", "mqtt", "*", "1.0")
	serviceDescriptor := cref.NewDescriptor("mqtt-listener", "service", "default", "*", "1.0")

	c.RegisterType(mqttMessageBusDescriptor, messagebus.NewMqttMessageBus)
	c.RegisterType(serviceDescriptor, service.NewMqttListenerService)

	return &c
}
