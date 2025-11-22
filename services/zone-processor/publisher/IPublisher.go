package publisher

import (
	"context"
)

type IPublisher interface {
	SendEvent(ctx context.Context, event any, msgType string) error
}
