package operations1

import (
	"context"
	"net/http"

	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/geo-history/clients/version1"
	"github.com/Shuv1Wolf/subterra-locate/facades/geo-facade/utils"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	httpcontr "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type GeoHistoryOperationsV1 struct {
	*httpcontr.RestOperations
	beaconAdmin clients1.IGeoHistoryClientV1
}

func NewGeoHistoryOperationsV1() *GeoHistoryOperationsV1 {
	c := GeoHistoryOperationsV1{
		RestOperations: httpcontr.NewRestOperations(),
	}
	c.DependencyResolver.Put(context.Background(), "geo-history", cref.NewDescriptor("geo-history", "client", "*", "*", "1.0"))
	return &c
}

func (c *GeoHistoryOperationsV1) SetReferences(ctx context.Context, references cref.IReferences) {
	c.RestOperations.SetReferences(ctx, references)

	dependency, _ := c.DependencyResolver.GetOneRequired("geo-history")
	client, ok := dependency.(clients1.IGeoHistoryClientV1)
	if !ok {
		panic("GeoHistoryOperationsV1: Cant't resolv dependency 'client' to IGeoHistoryClientV1")
	}
	c.beaconAdmin = client
}

func (c *GeoHistoryOperationsV1) GetHistory(res http.ResponseWriter, req *http.Request) {
	var paging = c.GetPagingParams(req)
	var sort = utils.GetSortFieldParams(req)

	reqctx := cdata.GetRequestContextParams(req)
	map_id := c.GetParam(req, "map_id")
	from := c.GetParam(req, "from")
	to := c.GetParam(req, "to")

	page, err := c.beaconAdmin.GetHistory(
		context.Background(), *reqctx, map_id, from, to, paging, sort,
	)

	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, page, nil)
	}
}
