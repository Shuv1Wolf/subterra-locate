package clients1

import (
	"context"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-history/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type GeoHistoryNullClientV1 struct {
}

func NewGeoHistoryNullClientV1() *GeoHistoryNullClientV1 {
	return &GeoHistoryNullClientV1{}
}

func (c *GeoHistoryNullClientV1) GetHistory(ctx context.Context, reqctx cdata.RequestContextV1, mapId, entityId, from, to string,
	paging *cquery.PagingParams, sortField *cquery.SortField) (cquery.DataPage[data1.HistoricalRecordV1], error) {
	return *cquery.NewEmptyDataPage[data1.HistoricalRecordV1](), nil
}
