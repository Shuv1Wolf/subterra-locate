package build

import (
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/generator/publisher"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/generator/service"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type GenEngineServiceFactory struct {
	*cbuild.Factory
}

func NewGenEngineServiceFactory() *GenEngineServiceFactory {
	c := GenEngineServiceFactory{}
	c.Factory = cbuild.NewFactory()

	natsPublisherDescriptor := cref.NewDescriptor("generator", "publisher", "nats", "*", "1.0")
	serviceDescriptor := cref.NewDescriptor("generator", "service", "default", "*", "1.0")

	c.RegisterType(natsPublisherDescriptor, publisher.NewNatsPublisher)
	c.RegisterType(serviceDescriptor, service.NewGenService)

	return &c
}
