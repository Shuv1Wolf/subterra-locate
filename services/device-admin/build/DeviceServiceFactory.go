package build

import (
	controllers "github.com/Shuv1Wolf/subterra-locate/services/device-admin/controllers/version1"
	persist "github.com/Shuv1Wolf/subterra-locate/services/device-admin/persistence"
	"github.com/Shuv1Wolf/subterra-locate/services/device-admin/publisher"
	logic "github.com/Shuv1Wolf/subterra-locate/services/device-admin/service"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type DeviceServiceFactory struct {
	cbuild.Factory
}

func NewDeviceServiceFactory() *DeviceServiceFactory {
	c := &DeviceServiceFactory{
		Factory: *cbuild.NewFactory(),
	}

	memoryPersistenceDescriptor := cref.NewDescriptor("device-admin", "persistence", "memory", "*", "1.0")
	postgresPersistenceDescriptor := cref.NewDescriptor("device-admin", "persistence", "postgres", "*", "1.0")
	serviceDescriptor := cref.NewDescriptor("device-admin", "service", "default", "*", "1.0")
	httpcontrollerV1Descriptor := cref.NewDescriptor("device-admin", "controller", "grpc", "*", "1.0")
	deviceEventsPuvlisher := cref.NewDescriptor("device-admin", "publisher", "nats", "device-events", "1.0")

	c.RegisterType(postgresPersistenceDescriptor, persist.NewDevicePostgresPersistence)
	c.RegisterType(memoryPersistenceDescriptor, persist.NewDeviceMemoryPersistence)
	c.RegisterType(serviceDescriptor, logic.NewDeviceService)
	c.RegisterType(httpcontrollerV1Descriptor, controllers.NewDeviceCommandableGrpcControllerV1)
	c.RegisterType(deviceEventsPuvlisher, publisher.NewNatsPublisher)

	return c
}
