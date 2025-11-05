package clients1

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-history/data/version1"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
)

type GeoHistoryGrpcClientV1 struct {
	*cclients.CommandableGrpcClient
}

func NewGeoHistoryGrpcClientV1() *GeoHistoryGrpcClientV1 {
	c := &GeoHistoryGrpcClientV1{
		CommandableGrpcClient: cclients.NewCommandableGrpcClient("geo.history.v1"),
	}
	return c
}

func (c *GeoHistoryGrpcClientV1) GetHistory(ctx context.Context, orgId, mapId, from, to string,
	paging *cquery.PagingParams, sortField *cquery.SortField) (cquery.DataPage[data1.HistoricalRecordV1], error) {

	var pagingMap map[string]interface{}
	if paging != nil {
		pagingMap = map[string]interface{}{
			"skip":  paging.Skip,
			"take":  paging.Take,
			"total": paging.Total,
		}
	}

	var sortMap map[string]interface{}
	if sortField != nil {
		sortMap = map[string]interface{}{
			"name":      sortField.Name,
			"ascending": sortField.Ascending,
		}
	}

	params := cdata.NewAnyValueMapFromTuples(
		"org_id", orgId,
		"map_id", mapId,
		"from", from,
		"to", to,
		"sort", sortMap,
		"paging", pagingMap,
	)

	response, err := c.CallCommand(ctx, "get_device_history", cdata.NewAnyValueMapFromValue(params.Value()))

	if err != nil {
		return *cquery.NewEmptyDataPage[data1.HistoricalRecordV1](), err
	}

	return cclients.HandleHttpResponse[cquery.DataPage[data1.HistoricalRecordV1]](response)
}
