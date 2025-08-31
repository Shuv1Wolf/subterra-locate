package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/generator/publisher"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

type GenService struct {
	pub    publisher.IPublisher
	isOpen bool
	cancel context.CancelFunc
	logger *log.ConsoleLogger
}

func NewGenService() *GenService {
	c := &GenService{}
	c.logger = log.NewConsoleLogger()
	return c
}

func (c *GenService) SetReferences(ctx context.Context, references cref.IReferences) {
	res, err := references.GetOneRequired(
		cref.NewDescriptor("generator", "publisher", "nats", "*", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.pub = res.(publisher.IPublisher)
}

func (c *GenService) Open(ctx context.Context) error {
	if c.isOpen {
		return nil
	}

	c.isOpen = true
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	go func() {
		for c.isOpen {
			delay := time.Duration(rand.Intn(3001)) * time.Millisecond
			time.Sleep(delay)

			if !c.isOpen {
				break
			}

			err := c.pub.SendPosition(ctx)
			if err != nil {
				c.logger.Error(ctx, err, "Failed to send position")
			} else {
				c.logger.Trace(ctx, "Sent position after %v delay", delay)
			}
		}
	}()

	return nil
}

func (c *GenService) Close(ctx context.Context) error {
	if !c.isOpen {
		return nil
	}

	c.isOpen = false
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

func (c *GenService) IsOpen() bool {
	return c.isOpen
}

func (c *GenService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	// No configuration needed for now
}
