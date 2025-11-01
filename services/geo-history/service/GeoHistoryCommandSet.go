package service

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-history/data/version1"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type GeoHistoryCommandSet struct {
	ccmd.CommandSet
	controller   IGeoHistoryService
	mapConvertor cconv.IJSONEngine[data1.HistoricalRecordV1]
}

func NewGeoHistoryCommandSet(controller IGeoHistoryService) *GeoHistoryCommandSet {
	c := &GeoHistoryCommandSet{
		CommandSet:   *ccmd.NewCommandSet(),
		controller:   controller,
		mapConvertor: cconv.NewDefaultCustomTypeJsonConvertor[data1.HistoricalRecordV1](),
	}

	c.AddCommand(c.makeGetHistoryCommand())

	return c
}

func (c *GeoHistoryCommandSet) makeGetHistoryCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"get_maps",
		cvalid.NewObjectSchema().
			WithRequiredProperty("org_id", cconv.String).
			WithRequiredProperty("map_id", cconv.String).
			WithRequiredProperty("from", cconv.String).
			WithRequiredProperty("to", cconv.String).
			WithOptionalProperty("sort", data1.NewSortFieldSchema()).
			WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			paging := cquery.NewEmptyPagingParams()
			if _val, ok := args.Get("paging"); ok {
				paging = cquery.NewPagingParamsFromValue(_val)
			}
			sort := cquery.NewEmptySortField()
			if _val, ok := args.Get("sort"); ok {
				sort = *data1.NewSortFieldFromValue(_val)
			}
			return c.controller.GetHistory(
				ctx,
				args.GetAsString("org_id"),
				args.GetAsString("map_id"),
				args.GetAsString("from"),
				args.GetAsString("to"),
				*paging,
				sort,
			)
		})
}
