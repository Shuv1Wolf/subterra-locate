package publisher

import (
	"context"
)

type IPublisher interface {
	SendBeaconCreatedEvent(ctx context.Context, id string) error
	SendBeaconChangedEvent(ctx context.Context, id string) error
	SendBeaconDeletedEvent(ctx context.Context, id string) error
}
