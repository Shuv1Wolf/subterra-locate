package controllers

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	grpcctrl "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/controllers"
)

type BeaconCommandableGrpcControllerV1 struct {
	*grpcctrl.CommandableGrpcController
}

func NewBeaconCommandableGrpcControllerV1() *BeaconCommandableGrpcControllerV1 {
	c := &BeaconCommandableGrpcControllerV1{}
	c.CommandableGrpcController = grpcctrl.InheritCommandableGrpcController(c, "beacon.admin.v1")
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("beacon-admin", "service", "*", "*", "1.0"))
	return c
}
