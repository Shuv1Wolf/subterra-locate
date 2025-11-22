package service

import (
	"context"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	exec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type ZoneCommandSet struct {
	ccmd.CommandSet
	controller    IZoneService
	zoneConvertor cconv.IJSONEngine[*data1.ZoneV1]
}

func NewZoneCommandSet(controller IZoneService) *ZoneCommandSet {
	c := &ZoneCommandSet{
		CommandSet:    *ccmd.NewCommandSet(),
		controller:    controller,
		zoneConvertor: cconv.NewDefaultCustomTypeJsonConvertor[*data1.ZoneV1](),
	}

	c.AddCommand(c.makeGetZonesCommand())
	c.AddCommand(c.makeGetZoneByIdCommand())
	c.AddCommand(c.makeCreateZoneCommand())
	c.AddCommand(c.makeUpdateZoneCommand())
	c.AddCommand(c.makeDeleteZoneByIdCommand())

	return c
}

func (c *ZoneCommandSet) makeGetZonesCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"get_zones",
		cvalid.NewObjectSchema().
			WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).
			WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			filter := cquery.NewEmptyFilterParams()
			paging := cquery.NewEmptyPagingParams()
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("filter"); ok {
				filter = cquery.NewFilterParamsFromValue(_val)
			}
			if _val, ok := args.Get("paging"); ok {
				paging = cquery.NewPagingParamsFromValue(_val)
			}
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.GetZones(ctx, *reqctx, *filter, *paging)
		})
}

func (c *ZoneCommandSet) makeGetZoneByIdCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"get_zone_by_id",
		cvalid.NewObjectSchema().
			WithRequiredProperty("zone_id", cconv.String).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.GetZoneById(ctx, *reqctx, args.GetAsString("zone_id"))
		})
}

func (c *ZoneCommandSet) makeCreateZoneCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"create_zone",
		cvalid.NewObjectSchema().
			WithRequiredProperty("zone", data1.NewZoneV1Schema()).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			var zone *data1.ZoneV1
			if _zone, ok := args.GetAsObject("zone"); ok {
				buf, err := cconv.JsonConverter.ToJson(_zone)
				if err != nil {
					return nil, err
				}
				zone, err = c.zoneConvertor.FromJson(buf)
				if err != nil {
					return nil, err
				}
			}
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.CreateZone(ctx, *reqctx, *zone)
		})
}

func (c *ZoneCommandSet) makeUpdateZoneCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"update_zone",
		cvalid.NewObjectSchema().
			WithRequiredProperty("zone", data1.NewZoneV1Schema()).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			var zone *data1.ZoneV1
			if _zone, ok := args.GetAsObject("zone"); ok {
				buf, err := cconv.JsonConverter.ToJson(_zone)
				if err != nil {
					return nil, err
				}
				zone, err = c.zoneConvertor.FromJson(buf)
				if err != nil {
					return nil, err
				}
			}
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.UpdateZone(ctx, *reqctx, *zone)
		})
}

func (c *ZoneCommandSet) makeDeleteZoneByIdCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"delete_zone_by_id",
		cvalid.NewObjectSchema().
			WithRequiredProperty("zone_id", cconv.String).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.DeleteZoneById(ctx, *reqctx, args.GetAsString("zone_id"))
		})
}
