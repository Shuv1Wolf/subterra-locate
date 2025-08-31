package publisher

import (
	"context"
)

type IPublisher interface {
	SendPosition(ctx context.Context) error
}
