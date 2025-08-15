package publisher

import (
	"context"

	nats "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
)

type IPublisher interface {
	SendDevicePosition(ctx context.Context, event *nats.DevicePositioningEventV1) error
}
