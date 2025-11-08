package operations1

import (
	"context"
	"net/http"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	httpcontr "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"

	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/beacon-admin/clients/version1"
	services1 "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
)

type BeaconAdminOperationsV1 struct {
	*httpcontr.RestOperations
	beaconAdmin clients1.IBeaconsClientV1
}

func NewBeaconAdminOperationsV1() *BeaconAdminOperationsV1 {
	c := BeaconAdminOperationsV1{
		RestOperations: httpcontr.NewRestOperations(),
	}
	c.DependencyResolver.Put(context.Background(), "beacon-admin", cref.NewDescriptor("beacon-admin", "client", "*", "*", "1.0"))
	return &c
}

func (c *BeaconAdminOperationsV1) SetReferences(ctx context.Context, references cref.IReferences) {
	c.RestOperations.SetReferences(ctx, references)

	dependency, _ := c.DependencyResolver.GetOneRequired("beacon-admin")
	client, ok := dependency.(clients1.IBeaconsClientV1)
	if !ok {
		panic("BeaconAdminOperationsV1: Cant't resolv dependency 'client' to IBeaconsClientV1")
	}
	c.beaconAdmin = client
}

func (c *BeaconAdminOperationsV1) GetBeacons(res http.ResponseWriter, req *http.Request) {
	var filter = c.GetFilterParams(req)
	var paging = c.GetPagingParams(req)
	var reqctx = cdata.GetRequestContextParams(req)

	page, err := c.beaconAdmin.GetBeacons(
		context.Background(), *reqctx, filter, paging)

	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, page, nil)
	}
}

func (c *BeaconAdminOperationsV1) GetBeaconById(res http.ResponseWriter, req *http.Request) {
	id := c.GetParam(req, "id")
	reqctx := cdata.GetRequestContextParams(req)

	item, err := c.beaconAdmin.GetBeaconById(context.Background(), *reqctx, id)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *BeaconAdminOperationsV1) GetBeaconByUdi(res http.ResponseWriter, req *http.Request) {
	udi := c.GetParam(req, "udi")
	reqctx := cdata.GetRequestContextParams(req)

	item, err := c.beaconAdmin.GetBeaconByUdi(context.Background(), *reqctx, udi)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *BeaconAdminOperationsV1) CreateBeacon(res http.ResponseWriter, req *http.Request) {
	reqctx := cdata.GetRequestContextParams(req)

	data := services1.BeaconV1{}
	err := c.DecodeBody(req, &data)
	if err != nil {
		c.SendError(res, req, err)
	}
	item, err := c.beaconAdmin.CreateBeacon(context.Background(), *reqctx, data)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *BeaconAdminOperationsV1) UpdateBeacon(res http.ResponseWriter, req *http.Request) {
	reqctx := cdata.GetRequestContextParams(req)

	data := services1.BeaconV1{}
	err := c.DecodeBody(req, &data)
	if err != nil {
		c.SendError(res, req, err)
	}

	item, err := c.beaconAdmin.UpdateBeacon(context.Background(), *reqctx, data)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *BeaconAdminOperationsV1) DeleteBeaconById(res http.ResponseWriter, req *http.Request) {
	id := c.GetParam(req, "id")
	reqctx := cdata.GetRequestContextParams(req)

	item, err := c.beaconAdmin.DeleteBeaconById(context.Background(), *reqctx, id)

	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}
