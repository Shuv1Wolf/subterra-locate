package controllers1

import (
	"context"
	"net/http"

	operations1 "github.com/Shuv1Wolf/subterra-locate/facades/geo-facade/operations/version1"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	httpcontr "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type FacadeControllerV1 struct {
	*httpcontr.RestController
	beaconsOperations *operations1.BeaconAdminOperationsV1
}

func NewFacadeControllerV1() *FacadeControllerV1 {
	c := &FacadeControllerV1{
		beaconsOperations: operations1.NewBeaconAdminOperationsV1(),
	}
	c.RestController = httpcontr.InheritRestController(c)
	c.BaseRoute = "api/v1/geo"
	return c
}

func (c *FacadeControllerV1) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.RestController.Configure(ctx, config)

	c.beaconsOperations.Configure(ctx, config)
}

func (c *FacadeControllerV1) SetReferences(ctx context.Context, references cref.IReferences) {
	c.RestController.SetReferences(ctx, references)

	c.beaconsOperations.SetReferences(ctx, references)
}

func (c *FacadeControllerV1) Register() {
	c.FacadeControllerV1()
}

func (c *FacadeControllerV1) FacadeControllerV1() {
	// Beacons routes
	c.RegisterRoute("get", "/beacons", nil,
		func(res http.ResponseWriter, req *http.Request) { c.beaconsOperations.GetBeacons(res, req) })
	c.RegisterRoute("get", "/beacons/:id", nil,
		func(res http.ResponseWriter, req *http.Request) { c.beaconsOperations.GetBeaconById(res, req) })
	c.RegisterRoute("get", "/beacons/udi/:udi", nil,
		func(res http.ResponseWriter, req *http.Request) { c.beaconsOperations.GetBeaconByUdi(res, req) })
	c.RegisterRoute("post", "/beacons", nil,
		func(res http.ResponseWriter, req *http.Request) { c.beaconsOperations.CreateBeacon(res, req) })
	c.RegisterRoute("put", "/beacons", nil,
		func(res http.ResponseWriter, req *http.Request) { c.beaconsOperations.UpdateBeacon(res, req) })
	c.RegisterRoute("delete", "/beacons/:id", nil,
		func(res http.ResponseWriter, req *http.Request) { c.beaconsOperations.DeleteBeaconById(res, req) })
}
