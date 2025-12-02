package listener

import (
	nqueues "github.com/pip-services4/pip-services4-go/pip-services4-nats-go/queues"
)

type NatsListener struct {
	*nqueues.NatsMessageQueue
}

func NewNatsListener() *NatsListener {
	c := &NatsListener{}
	c.NatsMessageQueue = nqueues.NewNatsMessageQueue("")
	return c
}
