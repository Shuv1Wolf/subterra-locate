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
	beaconsOperations     *operations1.BeaconAdminOperationsV1
	locationOperations    *operations1.LocationEngineOperationsV1
	geoRendererOperations *operations1.GeoRendererOperationsV1
	geoHistoryOperations  *operations1.GeoHistoryOperationsV1
	zoneOperations        *operations1.ZoneProcessorOperationsV1
	zoneMonitor           *operations1.ZoneMonitorOperationsV1
}

func NewFacadeControllerV1() *FacadeControllerV1 {
	c := &FacadeControllerV1{
		beaconsOperations:     operations1.NewBeaconAdminOperationsV1(),
		locationOperations:    operations1.NewLocationEngineOperationsV1(),
		geoRendererOperations: operations1.NewGeoRendererOperationsV1(),
		geoHistoryOperations:  operations1.NewGeoHistoryOperationsV1(),
		zoneOperations:        operations1.NewZoneProcessorOperationsV1(),
		zoneMonitor:           operations1.NewZoneMonitorOperationsV1(),
	}
	c.RestController = httpcontr.InheritRestController(c)
	c.BaseRoute = "api/v1/geo"
	return c
}

func (c *FacadeControllerV1) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.RestController.Configure(ctx, config)

	c.beaconsOperations.Configure(ctx, config)
	c.locationOperations.Configure(ctx, config)
	c.geoRendererOperations.Configure(ctx, config)
	c.geoHistoryOperations.Configure(ctx, config)
	c.zoneOperations.Configure(ctx, config)
	c.zoneMonitor.Configure(ctx, config)
}

func (c *FacadeControllerV1) SetReferences(ctx context.Context, references cref.IReferences) {
	c.RestController.SetReferences(ctx, references)

	c.beaconsOperations.SetReferences(ctx, references)
	c.locationOperations.SetReferences(ctx, references)
	c.geoRendererOperations.SetReferences(ctx, references)
	c.geoHistoryOperations.SetReferences(ctx, references)
	c.zoneOperations.SetReferences(ctx, references)
	c.zoneMonitor.SetReferences(ctx, references)
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

	// Locations routes
	c.RegisterRoute("get", "/location/device/monitor", nil,
		func(res http.ResponseWriter, req *http.Request) {
			c.locationOperations.MonitorDeviceLocationWS(res, req)
		})

	c.RegisterRoute("get", "/location/beacon/monitor", nil,
		func(res http.ResponseWriter, req *http.Request) {
			c.locationOperations.MonitorBeaconLocationWS(res, req)
		})

	// Geo renderer routes
	c.RegisterRoute("get", "/map", nil,
		func(res http.ResponseWriter, req *http.Request) { c.geoRendererOperations.GetMaps(res, req) })
	c.RegisterRoute("get", "/map/:id", nil,
		func(res http.ResponseWriter, req *http.Request) { c.geoRendererOperations.GetMapById(res, req) })
	c.RegisterRoute("post", "/map", nil,
		func(res http.ResponseWriter, req *http.Request) { c.geoRendererOperations.CreateMap(res, req) })
	c.RegisterRoute("post", "/map/upload", nil,
		func(res http.ResponseWriter, req *http.Request) { c.geoRendererOperations.UploadMapSVG(res, req) })
	c.RegisterRoute("put", "/map", nil,
		func(res http.ResponseWriter, req *http.Request) { c.geoRendererOperations.UpdateMap(res, req) })
	c.RegisterRoute("delete", "/map/:id", nil,
		func(res http.ResponseWriter, req *http.Request) { c.geoRendererOperations.DeleteMapById(res, req) })

	// Geo history routes
	c.RegisterRoute("get", "/history", nil,
		func(res http.ResponseWriter, req *http.Request) { c.geoHistoryOperations.GetHistory(res, req) })

	// Zone routes
	c.RegisterRoute("get", "/zones/monitor", nil,
		func(res http.ResponseWriter, req *http.Request) {
			c.zoneMonitor.MonitorZoneWS(res, req)
		})
	c.RegisterRoute("get", "/zones", nil,
		func(res http.ResponseWriter, req *http.Request) { c.zoneOperations.GetZones(res, req) })
	c.RegisterRoute("get", "/zones/:id", nil,
		func(res http.ResponseWriter, req *http.Request) { c.zoneOperations.GetZoneById(res, req) })
	c.RegisterRoute("post", "/zones", nil,
		func(res http.ResponseWriter, req *http.Request) { c.zoneOperations.CreateZone(res, req) })
	c.RegisterRoute("put", "/zones", nil,
		func(res http.ResponseWriter, req *http.Request) { c.zoneOperations.UpdateZone(res, req) })
	c.RegisterRoute("delete", "/zones/:id", nil,
		func(res http.ResponseWriter, req *http.Request) { c.zoneOperations.DeleteZoneById(res, req) })
}
