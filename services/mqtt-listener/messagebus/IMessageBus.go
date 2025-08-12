package messagebus

import (
	"context"

	nats "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
)

type IMessageBus interface {
	SendRawBle(ctx context.Context, event *nats.BLEBeaconRawEventV1) error
}
