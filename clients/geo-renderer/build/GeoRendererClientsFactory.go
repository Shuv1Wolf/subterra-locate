package build

import (
	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/geo-renderer/clients/version1"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type GeoRendererClientsFactory struct {
	cbuild.Factory
	NullClientDescriptor   *cref.Descriptor
	GrpcClientDescriptor   *cref.Descriptor
	MemoryClientDescriptor *cref.Descriptor
}

func NewGeoRendererClientsFactory() *GeoRendererClientsFactory {

	bcf := GeoRendererClientsFactory{}
	bcf.Factory = *cbuild.NewFactory()

	bcf.NullClientDescriptor = cref.NewDescriptor("geo-renderer", "client", "null", "*", "1.0")
	bcf.GrpcClientDescriptor = cref.NewDescriptor("geo-renderer", "client", "grpc", "*", "1.0")
	bcf.MemoryClientDescriptor = cref.NewDescriptor("geo-renderer", "client", "memory", "*", "1.0")

	bcf.RegisterType(bcf.NullClientDescriptor, clients1.NewGeoRendererNullClientV1)
	bcf.RegisterType(bcf.GrpcClientDescriptor, clients1.NewGeoRendererGrpcClientV1)
	bcf.RegisterType(bcf.MemoryClientDescriptor, clients1.NewGeoRendererMemoryClientV1)

	return &bcf
}
