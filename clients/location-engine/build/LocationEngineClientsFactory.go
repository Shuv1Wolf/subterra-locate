package build

import (
	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/location-engine/clients/version1"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type LocationEngineClientsFactory struct {
	cbuild.Factory
	GrpcClientDescriptor *cref.Descriptor
}

func NewLocationEngineClientsFactory() *LocationEngineClientsFactory {

	bcf := LocationEngineClientsFactory{}
	bcf.Factory = *cbuild.NewFactory()

	bcf.GrpcClientDescriptor = cref.NewDescriptor("location-monitor", "client", "grpc", "*", "1.0")

	bcf.RegisterType(bcf.GrpcClientDescriptor, clients1.NewLocationMonitorGrpcClientV1)

	return &bcf
}
