package build

import (
	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/device-admin/clients/version1"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type DeviceClientFactory struct {
	cbuild.Factory
	NullClientDescriptor   *cref.Descriptor
	GrpcClientDescriptor   *cref.Descriptor
	MemoryClientDescriptor *cref.Descriptor
}

func NewDeviceClientFactory() *DeviceClientFactory {

	bcf := DeviceClientFactory{}
	bcf.Factory = *cbuild.NewFactory()

	bcf.NullClientDescriptor = cref.NewDescriptor("device-admin", "client", "null", "*", "1.0")
	bcf.GrpcClientDescriptor = cref.NewDescriptor("device-admin", "client", "grpc", "*", "1.0")
	bcf.MemoryClientDescriptor = cref.NewDescriptor("device-admin", "client", "memory", "*", "1.0")

	bcf.RegisterType(bcf.NullClientDescriptor, clients1.NewDeviceNullClientV1)
	bcf.RegisterType(bcf.GrpcClientDescriptor, clients1.NewDeviceGrpcClientV1)
	bcf.RegisterType(bcf.MemoryClientDescriptor, clients1.NewDeviceMemoryClientV1)

	return &bcf
}
