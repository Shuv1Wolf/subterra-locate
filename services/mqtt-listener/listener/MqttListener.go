package listener

import (
	mqueues "github.com/pip-services4/pip-services4-go/pip-services4-mqtt-go/queues"
)

type MqttListener struct {
	*mqueues.MqttMessageQueue
}

func NewMqttListener() *MqttListener {
	c := &MqttListener{}
	c.MqttMessageQueue = mqueues.NewMqttMessageQueue("")
	return c
}
