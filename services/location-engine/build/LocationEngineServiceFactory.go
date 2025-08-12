package build

import (
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/listener"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/publisher"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/service"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type LocationEngineServiceFactory struct {
	*cbuild.Factory
}

func NewLocationEngineServiceFactory() *LocationEngineServiceFactory {
	c := LocationEngineServiceFactory{}
	c.Factory = cbuild.NewFactory()

	natsRawBleListenerDescriptor := cref.NewDescriptor("location-engine", "listener", "nats", "loc-raw-ble", "1.0")
	natsHistoryBlePublisherDescriptor := cref.NewDescriptor("location-engine", "publisher", "nats", "loc-hist-ble", "1.0")
	serviceDescriptor := cref.NewDescriptor("location-engine", "service", "default", "*", "1.0")

	c.RegisterType(natsRawBleListenerDescriptor, listener.NewNatsListener)
	c.RegisterType(serviceDescriptor, service.NewLocationEngineService)
	c.RegisterType(natsHistoryBlePublisherDescriptor, publisher.NewNatsPublisher)

	return &c
}
