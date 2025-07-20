package messagebus

import (
	"context"

	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
)

type IMessageBus interface {
	Listen(ctx context.Context, receiver cqueues.IMessageReceiver) error
	EndListen(ctx context.Context)
}
