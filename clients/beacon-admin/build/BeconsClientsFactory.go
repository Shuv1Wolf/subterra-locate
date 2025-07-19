package build

import (
	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/beacon-admin/clients/version1"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type BeaconsClientFactory struct {
	cbuild.Factory
	NullClientDescriptor   *cref.Descriptor
	DirectClientDescriptor *cref.Descriptor
	GrpcClientDescriptor   *cref.Descriptor
	MemoryClientDescriptor *cref.Descriptor
}

func NewBeaconsClientFactory() *BeaconsClientFactory {

	bcf := BeaconsClientFactory{}
	bcf.Factory = *cbuild.NewFactory()

	bcf.NullClientDescriptor = cref.NewDescriptor("beacon-admin", "client", "null", "*", "1.0")
	bcf.DirectClientDescriptor = cref.NewDescriptor("beacon-admin", "client", "direct", "*", "1.0")
	bcf.GrpcClientDescriptor = cref.NewDescriptor("beacon-admin", "client", "grpc", "*", "1.0")
	bcf.MemoryClientDescriptor = cref.NewDescriptor("beacon-admin", "client", "memory", "*", "1.0")

	bcf.RegisterType(bcf.NullClientDescriptor, clients1.NewBeaconsNullClientV1)
	bcf.RegisterType(bcf.DirectClientDescriptor, clients1.NewBeaconsDirectClientV1)
	bcf.RegisterType(bcf.GrpcClientDescriptor, clients1.NewBeaconsGrpcClientV1)
	bcf.RegisterType(bcf.MemoryClientDescriptor, clients1.NewBeaconsMemoryClientV1)

	return &bcf
}
