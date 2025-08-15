package publisher

import (
	"context"
	"encoding/json"
	"time"

	natsConst "github.com/Shuv1Wolf/subterra-locate/services/common/nats/const"
	natsEvents "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	nqueues "github.com/pip-services4/pip-services4-go/pip-services4-nats-go/queues"
)

type NatsPublisher struct {
	*nqueues.NatsMessageQueue
}

func NewNatsPublisher() *NatsPublisher {
	c := &NatsPublisher{}
	c.NatsMessageQueue = nqueues.NewNatsMessageQueue("")
	return c
}

func (c *NatsPublisher) SendDevicePosition(ctx context.Context, event *natsEvents.DevicePositioningEventV1) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to serialize message")
		return err
	}

	envelope := &cqueues.MessageEnvelope{
		MessageId:   keys.IdGenerator.NextShort(),
		SentTime:    time.Now(),
		TraceId:     cctx.GetTraceId(ctx),
		MessageType: natsConst.NATS_EVENTS_DEVICE_POSITION,
		Message:     bytes,
	}

	err = c.Send(ctx, envelope)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to send message")
		return err
	}

	return nil
}
