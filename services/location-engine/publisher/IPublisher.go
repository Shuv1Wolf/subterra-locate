package publisher

import (
	"context"

	nats "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
)

type IPublisher interface {
	SendHistoryBle(ctx context.Context, event *nats.BLEBeaconHistoryEventV1) error
}
