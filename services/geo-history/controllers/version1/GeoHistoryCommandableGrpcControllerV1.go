package controllers

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	grpcctrl "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/controllers"
)

type GeoHistoryCommandableGrpcControllerV1 struct {
	*grpcctrl.CommandableGrpcController
}

func NewGeoHistoryCommandableGrpcControllerV1() *GeoHistoryCommandableGrpcControllerV1 {
	c := &GeoHistoryCommandableGrpcControllerV1{}
	c.CommandableGrpcController = grpcctrl.InheritCommandableGrpcController(c, "geo.history.v1")
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("geo-history", "service", "*", "*", "1.0"))
	return c
}
