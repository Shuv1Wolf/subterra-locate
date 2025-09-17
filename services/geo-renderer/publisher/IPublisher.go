package publisher

import (
	"context"
)

type IPublisher interface {
	SendMap2dCreatedEvent(ctx context.Context, id string) error
	SendMap2dChangedEvent(ctx context.Context, id string) error
	SendMap2dDeletedEvent(ctx context.Context, id string) error
}
