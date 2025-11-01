package service

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-history/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IGeoHistoryService interface {
	GetHistory(ctx context.Context, orgId, mapId, from, to string, paging cquery.PagingParams, sortField cquery.SortField) (cquery.DataPage[data1.HistoricalRecordV1], error)
}
