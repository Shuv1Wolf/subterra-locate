package build

import (
	controllers "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/controllers/version1"
	persist "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/persistence"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/publisher"
	logic "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/service"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type MapServiceFactory struct {
	cbuild.Factory
}

func NewMapServiceFactory() *MapServiceFactory {
	c := &MapServiceFactory{
		Factory: *cbuild.NewFactory(),
	}

	postgresPersistenceDescriptor := cref.NewDescriptor("geo-renderer", "persistence", "postgres", "map-2d", "1.0")
	serviceDescriptor := cref.NewDescriptor("geo-renderer", "service", "default", "*", "1.0")
	grpcControllerV1Descriptor := cref.NewDescriptor("geo-renderer", "controller", "grpc", "*", "1.0")
	mapEventsPublisher := cref.NewDescriptor("geo-renderer", "publisher", "nats", "map-2d-events", "1.0")

	c.RegisterType(postgresPersistenceDescriptor, persist.NewMap2dPostgresPersistence)
	c.RegisterType(serviceDescriptor, logic.NewMapService)
	c.RegisterType(grpcControllerV1Descriptor, controllers.NewMapCommandableGrpcControllerV1)
	c.RegisterType(mapEventsPublisher, publisher.NewNatsPublisher)

	return c
}
