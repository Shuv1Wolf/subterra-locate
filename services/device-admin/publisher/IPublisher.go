package publisher

import (
	"context"
)

type IPublisher interface {
	SendDeviceCreatedEvent(ctx context.Context, id string) error
	SendDeviceChangedEvent(ctx context.Context, id string) error
	SendDeviceDeletedEvent(ctx context.Context, id string) error
}
