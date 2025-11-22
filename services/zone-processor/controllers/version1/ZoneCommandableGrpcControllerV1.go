package version1

import (
	"context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	grpcctrl "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/controllers"
)

type ZoneCommandableGrpcControllerV1 struct {
	*grpcctrl.CommandableGrpcController
}

func NewZoneCommandableGrpcControllerV1() *ZoneCommandableGrpcControllerV1 {
	c := &ZoneCommandableGrpcControllerV1{}
	c.CommandableGrpcController = grpcctrl.InheritCommandableGrpcController(c, "zone_processor.v1")
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("zone-processor", "service", "*", "*", "1.0"))
	return c
}
