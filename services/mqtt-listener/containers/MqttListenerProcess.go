package containers

import (
	"github.com/Shuv1Wolf/subterra-locate/services/mqtt-listener/build"
	cproc "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
)

type MqttListenerProcess struct {
	*cproc.ProcessContainer
}

func NewMqttListenerProcess() *MqttListenerProcess {
	c := MqttListenerProcess{}
	c.ProcessContainer = cproc.NewProcessContainer("service-mqtt-listener", "MQTT listener service")
	c.AddFactory(build.NewMqttListenerServiceFactory())

	return &c
}
