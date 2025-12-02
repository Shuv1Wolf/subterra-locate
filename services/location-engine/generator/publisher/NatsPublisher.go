package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
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
	c.NatsMessageQueue = nqueues.NewNatsMessageQueue(natsConst.NATS_EVENTS_BLE_RSSI_TOPIC)
	rand.Seed(time.Now().UnixNano())
	return c
}

func (c *NatsPublisher) SendPosition(ctx context.Context) error {
	beacons := make([]natsEvents.BLEBeaconRawV1, 3)
	for i := 0; i < 3; i++ {
		beacons[i] = natsEvents.BLEBeaconRawV1{
			Id:      "beacon$10" + fmt.Sprintf("%d", i), // beacon$100, beacon$101, beacon$102
			Rssi:    rand.Intn(101) - 100,               // RSSI
			Txpower: 0,
		}
	}

	event := natsEvents.DeviceDetectedBLERawEventV1{
		DeviceMAC: "00:00:00:00:00:00",
		DeviceId:  "device$10" + fmt.Sprintf("%d", rand.Intn(3)), // device$100, device$101, device$102
		Count:     3,
		Beacons:   beacons,
	}

	c.Logger.Trace(ctx, "Generated event with beacons: %v", beacons)

	bytes, err := json.Marshal(event)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to serialize message")
		return err
	}

	envelope := &cqueues.MessageEnvelope{
		MessageId:   keys.IdGenerator.NextShort(),
		SentTime:    time.Now(),
		TraceId:     cctx.GetTraceId(ctx),
		MessageType: natsConst.NATS_EVENTS_BLE_RSSI_TOPIC,
		Message:     bytes,
	}

	err = c.Send(ctx, envelope)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to send message")
		return err
	}

	return nil
}
