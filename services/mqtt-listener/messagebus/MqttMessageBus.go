package messagebus

import (
	mqttconst "github.com/Shuv1Wolf/subterra-locate/services/common/mqtt-const"
	mqueues "github.com/pip-services4/pip-services4-go/pip-services4-mqtt-go/queues"
)

type MqttMessageBus struct {
	*mqueues.MqttMessageQueue
}

func NewMqttMessageBus() *MqttMessageBus {
	c := &MqttMessageBus{}
	c.MqttMessageQueue = mqueues.NewMqttMessageQueue(mqttconst.MQTT_BEACON_TOPIC)
	return c
}
