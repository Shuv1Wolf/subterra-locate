package operations1

import (
	"context"
	"net/http"
	"time"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	httpcontr "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"

	"io/ioutil"

	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/geo-renderer/clients/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	services1 "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
)

type GeoRendererOperationsV1 struct {
	*httpcontr.RestOperations
	geoRenderer clients1.IGeoRendererClientV1
}

func NewGeoRendererOperationsV1() *GeoRendererOperationsV1 {
	c := GeoRendererOperationsV1{
		RestOperations: httpcontr.NewRestOperations(),
	}
	c.DependencyResolver.Put(context.Background(), "geo-renderer", cref.NewDescriptor("geo-renderer", "client", "*", "*", "1.0"))
	return &c
}

func (c *GeoRendererOperationsV1) SetReferences(ctx context.Context, references cref.IReferences) {
	c.RestOperations.SetReferences(ctx, references)

	dependency, _ := c.DependencyResolver.GetOneRequired("geo-renderer")
	client, ok := dependency.(clients1.IGeoRendererClientV1)
	if !ok {
		panic("GeoRendererOperationsV1: Cant't resolv dependency 'client' to IGeoRendererClientV1")
	}
	c.geoRenderer = client
}

func (c *GeoRendererOperationsV1) GetMaps(res http.ResponseWriter, req *http.Request) {
	var filter = c.GetFilterParams(req)
	var paging = c.GetPagingParams(req)

	reqctx := cdata.GetRequestContextParams(req)

	page, err := c.geoRenderer.GetMaps(
		context.Background(), *reqctx, filter, paging)

	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, page, nil)
	}
}

func (c *GeoRendererOperationsV1) GetMapById(res http.ResponseWriter, req *http.Request) {
	id := c.GetParam(req, "id")
	reqctx := cdata.GetRequestContextParams(req)

	item, err := c.geoRenderer.GetMapById(context.Background(), *reqctx, id)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *GeoRendererOperationsV1) UploadMapSVG(res http.ResponseWriter, req *http.Request) {
	id := c.GetParam(req, "id")
	reqctx := cdata.GetRequestContextParams(req)

	if err := req.ParseMultipartForm(10 << 20); err != nil {
		c.SendError(res, req, err)
		return
	}
	file, _, err := req.FormFile("file")
	if err != nil {
		c.SendError(res, req, err)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.SendError(res, req, err)
		return
	}
	svgString := string(data)
	item, err := c.geoRenderer.GetMapById(context.Background(), *reqctx, id)
	if err != nil {
		c.SendError(res, req, err)
		return
	}
	item.SVG = svgString
	updated, err := c.geoRenderer.UpdateMap(context.Background(), *reqctx, *item)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, updated, nil)
	}
}

func (c *GeoRendererOperationsV1) CreateMap(res http.ResponseWriter, req *http.Request) {
	reqctx := cdata.GetRequestContextParams(req)

	data := services1.Map2dV1{}
	err := c.DecodeBody(req, &data)
	if err != nil {
		c.SendError(res, req, err)
	}
	data.CreatedAt = time.Now()
	item, err := c.geoRenderer.CreateMap(context.Background(), *reqctx, data)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *GeoRendererOperationsV1) UpdateMap(res http.ResponseWriter, req *http.Request) {
	reqctx := cdata.GetRequestContextParams(req)

	data := services1.Map2dV1{}
	err := c.DecodeBody(req, &data)
	if err != nil {
		c.SendError(res, req, err)
	}

	item, err := c.geoRenderer.UpdateMap(context.Background(), *reqctx, data)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *GeoRendererOperationsV1) DeleteMapById(res http.ResponseWriter, req *http.Request) {
	id := c.GetParam(req, "id")
	reqctx := cdata.GetRequestContextParams(req)

	item, err := c.geoRenderer.DeleteMapById(context.Background(), *reqctx, id)

	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}
