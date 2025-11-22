package operations1

import (
	"context"
	"net/http"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	httpcontr "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"

	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/zone-processor/clients/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
)

type ZoneProcessorOperationsV1 struct {
	*httpcontr.RestOperations
	zoneClient clients1.IZoneClientV1
}

func NewZoneProcessorOperationsV1() *ZoneProcessorOperationsV1 {
	c := ZoneProcessorOperationsV1{
		RestOperations: httpcontr.NewRestOperations(),
	}
	c.DependencyResolver.Put(context.Background(), "zone-processor", cref.NewDescriptor("zone-processor", "client", "*", "*", "1.0"))
	return &c
}

func (c *ZoneProcessorOperationsV1) SetReferences(ctx context.Context, references cref.IReferences) {
	c.RestOperations.SetReferences(ctx, references)

	dependency, _ := c.DependencyResolver.GetOneRequired("zone-processor")
	client, ok := dependency.(clients1.IZoneClientV1)
	if !ok {
		panic("ZoneProcessorOperationsV1: Can't resolve dependency 'client' to IZoneClientV1")
	}
	c.zoneClient = client
}

func (c *ZoneProcessorOperationsV1) GetZones(res http.ResponseWriter, req *http.Request) {
	var filter = c.GetFilterParams(req)
	var paging = c.GetPagingParams(req)
	var reqctx = cdata.GetRequestContextParams(req)

	page, err := c.zoneClient.GetZones(
		context.Background(), *reqctx, filter, paging)

	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, page, nil)
	}
}

func (c *ZoneProcessorOperationsV1) GetZoneById(res http.ResponseWriter, req *http.Request) {
	id := c.GetParam(req, "id")
	reqctx := cdata.GetRequestContextParams(req)

	item, err := c.zoneClient.GetZoneById(context.Background(), *reqctx, id)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *ZoneProcessorOperationsV1) CreateZone(res http.ResponseWriter, req *http.Request) {
	reqctx := cdata.GetRequestContextParams(req)

	data := data1.ZoneV1{}
	err := c.DecodeBody(req, &data)
	if err != nil {
		c.SendError(res, req, err)
		return
	}
	item, err := c.zoneClient.CreateZone(context.Background(), *reqctx, data)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *ZoneProcessorOperationsV1) UpdateZone(res http.ResponseWriter, req *http.Request) {
	reqctx := cdata.GetRequestContextParams(req)

	data := data1.ZoneV1{}
	err := c.DecodeBody(req, &data)
	if err != nil {
		c.SendError(res, req, err)
		return
	}

	item, err := c.zoneClient.UpdateZone(context.Background(), *reqctx, data)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *ZoneProcessorOperationsV1) DeleteZoneById(res http.ResponseWriter, req *http.Request) {
	id := c.GetParam(req, "id")
	reqctx := cdata.GetRequestContextParams(req)

	item, err := c.zoneClient.DeleteZoneById(context.Background(), *reqctx, id)

	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}
