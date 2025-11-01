package build

import (
	controllers "github.com/Shuv1Wolf/subterra-locate/services/geo-history/controllers/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-history/listener"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-history/persistence"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-history/service"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type GeoHistoryServiceFactory struct {
	*cbuild.Factory
}

func NewGeoHistoryServiceFactory() *GeoHistoryServiceFactory {
	c := &GeoHistoryServiceFactory{
		Factory: cbuild.NewFactory(),
	}

	eventsListenerDescriptor := cref.NewDescriptor("geo-history", "listener", "nats", "device-position", "1.0")
	postgresPersistenceDescriptor := cref.NewDescriptor("geo-history", "persistence", "postgres", "device", "1.0")
	serviceDescriptor := cref.NewDescriptor("geo-history", "service", "default", "*", "1.0")
	grpcControllerV1Descriptor := cref.NewDescriptor("geo-history", "controller", "grpc", "*", "1.0")

	c.RegisterType(eventsListenerDescriptor, listener.NewNatsListener)
	c.RegisterType(postgresPersistenceDescriptor, persistence.NewGeoHistoryPostgresPersistence)
	c.RegisterType(serviceDescriptor, service.NewGeoHistoryService)
	c.RegisterType(grpcControllerV1Descriptor, controllers.NewGeoHistoryCommandableGrpcControllerV1)

	return c
}
