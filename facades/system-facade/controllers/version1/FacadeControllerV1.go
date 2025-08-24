package controllers1

import (
	"context"
	"net/http"

	operations1 "github.com/Shuv1Wolf/subterra-locate/facades/system-facade/operations/version1"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	httpcontr "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type FacadeControllerV1 struct {
	*httpcontr.RestController
	deviceOperations *operations1.DeviceAdminOperationsV1
}

func NewFacadeControllerV1() *FacadeControllerV1 {
	c := &FacadeControllerV1{
		deviceOperations: operations1.NewDeviceAdminOperationsV1(),
	}
	c.RestController = httpcontr.InheritRestController(c)
	c.BaseRoute = "api/v1/system"
	return c
}

func (c *FacadeControllerV1) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.RestController.Configure(ctx, config)

	c.deviceOperations.Configure(ctx, config)
}

func (c *FacadeControllerV1) SetReferences(ctx context.Context, references cref.IReferences) {
	c.RestController.SetReferences(ctx, references)

	c.deviceOperations.SetReferences(ctx, references)
}

func (c *FacadeControllerV1) Register() {
	c.FacadeControllerV1()
}

func (c *FacadeControllerV1) FacadeControllerV1() {
	// Device routes
	c.RegisterRoute("get", "/devices", nil,
		func(res http.ResponseWriter, req *http.Request) { c.deviceOperations.GetDevices(res, req) })
	c.RegisterRoute("get", "/device/:id", nil,
		func(res http.ResponseWriter, req *http.Request) { c.deviceOperations.GetDeviceById(res, req) })
	c.RegisterRoute("post", "/device", nil,
		func(res http.ResponseWriter, req *http.Request) { c.deviceOperations.CreateDevice(res, req) })
	c.RegisterRoute("put", "/device", nil,
		func(res http.ResponseWriter, req *http.Request) { c.deviceOperations.UpdateDevice(res, req) })
	c.RegisterRoute("delete", "/device/:id", nil,
		func(res http.ResponseWriter, req *http.Request) { c.deviceOperations.DeleteDeviceById(res, req) })
}
