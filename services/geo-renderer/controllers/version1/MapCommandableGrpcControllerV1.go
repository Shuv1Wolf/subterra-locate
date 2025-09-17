package controllers

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	grpcctrl "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/controllers"
)

type MapCommandableGrpcControllerV1 struct {
	*grpcctrl.CommandableGrpcController
}

func NewMapCommandableGrpcControllerV1() *MapCommandableGrpcControllerV1 {
	c := &MapCommandableGrpcControllerV1{}
	c.CommandableGrpcController = grpcctrl.InheritCommandableGrpcController(c, "geo.renderer.v1")
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("geo-renderer", "service", "*", "*", "1.0"))
	return c
}
