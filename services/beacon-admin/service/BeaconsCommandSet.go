package logic

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	exec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type BeaconsCommandSet struct {
	ccmd.CommandSet
	controller      IBeaconsService
	beaconConvertor cconv.IJSONEngine[data.BeaconV1]
}

func NewBeaconsCommandSet(controller IBeaconsService) *BeaconsCommandSet {
	c := &BeaconsCommandSet{
		CommandSet:      *ccmd.NewCommandSet(),
		controller:      controller,
		beaconConvertor: cconv.NewDefaultCustomTypeJsonConvertor[data.BeaconV1](),
	}

	c.AddCommand(c.makeGetBeaconsCommand())
	c.AddCommand(c.makeGetBeaconByIdCommand())
	c.AddCommand(c.makeGetBeaconByUdiCommand())
	c.AddCommand(c.makeCreateBeaconCommand())
	c.AddCommand(c.makeUpdateBeaconCommand())
	c.AddCommand(c.makeDeleteBeaconByIdCommand())

	return c
}

func (c *BeaconsCommandSet) makeGetBeaconsCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"get_beacons",
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
			return c.controller.GetBeacons(ctx, *reqctx, *filter, *paging)
		})
}

func (c *BeaconsCommandSet) makeGetBeaconByIdCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"get_beacon_by_id",
		cvalid.NewObjectSchema().
			WithRequiredProperty("beacon_id", cconv.String).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.GetBeaconById(ctx, *reqctx, args.GetAsString("beacon_id"))
		})
}

func (c *BeaconsCommandSet) makeGetBeaconByUdiCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"get_beacon_by_udi",
		cvalid.NewObjectSchema().
			WithRequiredProperty("udi", cconv.String).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.GetBeaconByUdi(ctx, *reqctx, args.GetAsString("udi"))
		})
}

func (c *BeaconsCommandSet) makeCreateBeaconCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"create_beacon",
		cvalid.NewObjectSchema().
			WithRequiredProperty("beacon", data.NewBeaconV1Schema()).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {

			var beacon data.BeaconV1
			if _beacon, ok := args.GetAsObject("beacon"); ok {
				buf, err := cconv.JsonConverter.ToJson(_beacon)
				if err != nil {
					return nil, err
				}
				beacon, err = c.beaconConvertor.FromJson(buf)
				if err != nil {
					return nil, err
				}
			}
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.CreateBeacon(ctx, *reqctx, beacon)
		})
}

func (c *BeaconsCommandSet) makeUpdateBeaconCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"update_beacon",
		cvalid.NewObjectSchema().
			WithRequiredProperty("beacon", data.NewBeaconV1Schema()).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			var beacon data.BeaconV1
			if _beacon, ok := args.GetAsObject("beacon"); ok {
				buf, err := cconv.JsonConverter.ToJson(_beacon)
				if err != nil {
					return nil, err
				}
				beacon, err = c.beaconConvertor.FromJson(buf)
				if err != nil {
					return nil, err
				}
			}

			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.UpdateBeacon(ctx, *reqctx, beacon)
		})
}

func (c *BeaconsCommandSet) makeDeleteBeaconByIdCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"delete_beacon_by_id",
		cvalid.NewObjectSchema().
			WithRequiredProperty("beacon_id", cconv.String).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.DeleteBeaconById(ctx, *reqctx, args.GetAsString("beacon_id"))
		})
}
