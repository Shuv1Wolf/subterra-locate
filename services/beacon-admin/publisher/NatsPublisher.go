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

func (c *NatsPublisher) SendBeaconDeletedEvent(ctx context.Context, id string) error {
	return c.sendBeaconEvent(ctx, id, natsConst.NATS_BEACONS_EVENTS_DELETED_TYPE)
}

func (c *NatsPublisher) SendBeaconChangedEvent(ctx context.Context, id string) error {
	return c.sendBeaconEvent(ctx, id, natsConst.NATS_BEACONS_EVENTS_CHANGED_TYPE)
}

func (c *NatsPublisher) SendBeaconCreatedEvent(ctx context.Context, id string) error {
	return c.sendBeaconEvent(ctx, id, natsConst.NATS_BEACONS_EVENTS_CREATED_TYPE)
}

func (c *NatsPublisher) sendBeaconEvent(ctx context.Context, id string, msgType string) error {
	event := &natsEvents.BeaconChangedEvent{
		Id: id,
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to serialize message")
		return err
	}

	envelope := &cqueues.MessageEnvelope{
		MessageId:   keys.IdGenerator.NextShort(),
		SentTime:    time.Now(),
		TraceId:     cctx.GetTraceId(ctx),
		MessageType: msgType,
		Message:     bytes,
	}

	err = c.Send(ctx, envelope)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to send message")
		return err
	}

	return nil
}
