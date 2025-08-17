package build

import (
	controllers "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/controllers/version1"
	persist "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/persistence"
	"github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/publisher"
	logic "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/service"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type BeaconsServiceFactory struct {
	cbuild.Factory
}

func NewBeaconsServiceFactory() *BeaconsServiceFactory {
	c := &BeaconsServiceFactory{
		Factory: *cbuild.NewFactory(),
	}

	memoryPersistenceDescriptor := cref.NewDescriptor("beacon-admin", "persistence", "memory", "*", "1.0")
	postgresPersistenceDescriptor := cref.NewDescriptor("beacon-admin", "persistence", "postgres", "*", "1.0")
	serviceDescriptor := cref.NewDescriptor("beacon-admin", "service", "default", "*", "1.0")
	httpcontrollerV1Descriptor := cref.NewDescriptor("beacon-admin", "controller", "grpc", "*", "1.0")
	beaconsEventsPuvlisher := cref.NewDescriptor("beacon-admin", "publisher", "nats", "beacons-events", "1.0")

	c.RegisterType(postgresPersistenceDescriptor, persist.NewBeaconsPostgresPersistence)
	c.RegisterType(memoryPersistenceDescriptor, persist.NewBeaconsMemoryPersistence)
	c.RegisterType(serviceDescriptor, logic.NewBeaconsService)
	c.RegisterType(httpcontrollerV1Descriptor, controllers.NewBeaconCommandableGrpcControllerV1)
	c.RegisterType(beaconsEventsPuvlisher, publisher.NewNatsPublisher)

	return c
}
