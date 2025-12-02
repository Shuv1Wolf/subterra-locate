package build

import (
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/controllers/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/listener"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/persistence"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/publisher"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/service"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type ZoneServiceFactory struct {
	cbuild.Factory
}

func NewZoneServiceFactory() *ZoneServiceFactory {
	c := &ZoneServiceFactory{
		Factory: *cbuild.NewFactory(),
	}

	eventsListenerDescriptor := cref.NewDescriptor("zone-processor", "listener", "nats", "device-position", "1.0")
	postgresPersistenceDescriptor := cref.NewDescriptor("zone-processor", "persistence", "postgres", "*", "1.0")
	serviceDescriptor := cref.NewDescriptor("zone-processor", "service", "default", "*", "1.0")
	grpcControllerV1Descriptor := cref.NewDescriptor("zone-processor", "controller", "grpc", "*", "1.0")
	zoneEventsPublisher := cref.NewDescriptor("zone-processor", "publisher", "nats", "zone-events", "1.0")

	c.RegisterType(eventsListenerDescriptor, listener.NewNatsListener)
	c.RegisterType(postgresPersistenceDescriptor, persistence.NewZonePostgresPersistence)
	c.RegisterType(serviceDescriptor, service.NewZoneService)
	c.RegisterType(grpcControllerV1Descriptor, version1.NewZoneCommandableGrpcControllerV1)
	c.RegisterType(zoneEventsPublisher, publisher.NewNatsPublisher)

	return c
}
