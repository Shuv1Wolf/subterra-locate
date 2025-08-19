package controllers

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	grpcctrl "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/controllers"
)

type DeviceCommandableGrpcControllerV1 struct {
	*grpcctrl.CommandableGrpcController
}

func NewDeviceCommandableGrpcControllerV1() *DeviceCommandableGrpcControllerV1 {
	c := &DeviceCommandableGrpcControllerV1{}
	c.CommandableGrpcController = grpcctrl.InheritCommandableGrpcController(c, "device.admin.v1")
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("device-admin", "service", "*", "*", "1.0"))
	return c
}
