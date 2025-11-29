package build

import (
	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/zone-processor/clients/version1"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type ZoneProcessorClientsFactory struct {
	cbuild.Factory
	GrpcClientDescriptor        *cref.Descriptor
	GrpcMonitorClientDescriptor *cref.Descriptor
}

func NewZoneProcessorClientsFactory() *ZoneProcessorClientsFactory {

	bcf := ZoneProcessorClientsFactory{}
	bcf.Factory = *cbuild.NewFactory()

	bcf.GrpcClientDescriptor = cref.NewDescriptor("zone-processor", "client", "grpc", "*", "1.0")
	bcf.GrpcMonitorClientDescriptor = cref.NewDescriptor("zone-monitor", "client", "grpc", "*", "1.0")

	bcf.RegisterType(bcf.GrpcClientDescriptor, clients1.NewZoneGrpcClientV1)
	bcf.RegisterType(bcf.GrpcMonitorClientDescriptor, clients1.NewZoneMonitorGrpcClientV1)

	return &bcf
}
