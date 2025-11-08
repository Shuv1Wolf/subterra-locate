package persistence

import (
	"context"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-history/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IGeoHistoryPersistence interface {
	InsertBatch(ctx context.Context, items []*data1.HistoricalRecordV1) error

	GetHistory(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams, sortField cquery.SortField) (cquery.DataPage[data1.HistoricalRecordV1], error)
}
