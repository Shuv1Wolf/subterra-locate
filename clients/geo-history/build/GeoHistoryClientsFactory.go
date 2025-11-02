package build

import (
	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/geo-history/clients/version1"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type GeoHistoryClientsFactory struct {
	cbuild.Factory
	NullClientDescriptor   *cref.Descriptor
	GrpcClientDescriptor   *cref.Descriptor
	MemoryClientDescriptor *cref.Descriptor
}

func NewGeoHistoryClientsFactory() *GeoHistoryClientsFactory {

	bcf := GeoHistoryClientsFactory{}
	bcf.Factory = *cbuild.NewFactory()

	bcf.NullClientDescriptor = cref.NewDescriptor("geo-history", "client", "null", "*", "1.0")
	bcf.GrpcClientDescriptor = cref.NewDescriptor("geo-history", "client", "grpc", "*", "1.0")

	bcf.RegisterType(bcf.NullClientDescriptor, clients1.NewGeoHistoryNullClientV1)
	bcf.RegisterType(bcf.GrpcClientDescriptor, clients1.NewGeoHistoryGrpcClientV1)

	return &bcf
}
