package logic

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	exec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type MapCommandSet struct {
	ccmd.CommandSet
	controller   IMapService
	mapConvertor cconv.IJSONEngine[data.Map2dV1]
}

func NewMapCommandSet(controller IMapService) *MapCommandSet {
	c := &MapCommandSet{
		CommandSet:   *ccmd.NewCommandSet(),
		controller:   controller,
		mapConvertor: cconv.NewDefaultCustomTypeJsonConvertor[data.Map2dV1](),
	}

	c.AddCommand(c.makeGetMapsCommand())
	c.AddCommand(c.makeGetMapByIdCommand())
	c.AddCommand(c.makeCreateMapCommand())
	c.AddCommand(c.makeUpdateMapCommand())
	c.AddCommand(c.makeDeleteMapByIdCommand())

	return c
}

func (c *MapCommandSet) makeGetMapsCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"get_maps",
		cvalid.NewObjectSchema().
			WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).
			WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			filter := cquery.NewEmptyFilterParams()
			paging := cquery.NewEmptyPagingParams()
			if _val, ok := args.Get("filter"); ok {
				filter = cquery.NewFilterParamsFromValue(_val)
			}
			if _val, ok := args.Get("paging"); ok {
				paging = cquery.NewPagingParamsFromValue(_val)
			}
			return c.controller.GetMaps(ctx, *filter, *paging)
		})
}

func (c *MapCommandSet) makeGetMapByIdCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"get_map_by_id",
		cvalid.NewObjectSchema().
			WithRequiredProperty("map_id", cconv.String),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			return c.controller.GetMapById(ctx, args.GetAsString("map_id"))
		})
}

func (c *MapCommandSet) makeCreateMapCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"create_map",
		cvalid.NewObjectSchema().
			WithRequiredProperty("map", data.NewMap2dV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {

			var map2d data.Map2dV1
			if _map, ok := args.GetAsObject("map"); ok {
				buf, err := cconv.JsonConverter.ToJson(_map)
				if err != nil {
					return nil, err
				}
				map2d, err = c.mapConvertor.FromJson(buf)
				if err != nil {
					return nil, err
				}
			}
			return c.controller.CreateMap(ctx, map2d)
		})
}

func (c *MapCommandSet) makeUpdateMapCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"update_map",
		cvalid.NewObjectSchema().
			WithRequiredProperty("map", data.NewMap2dV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			var map2d data.Map2dV1
			if _map, ok := args.GetAsObject("map"); ok {
				buf, err := cconv.JsonConverter.ToJson(_map)
				if err != nil {
					return nil, err
				}
				map2d, err = c.mapConvertor.FromJson(buf)
				if err != nil {
					return nil, err
				}
			}
			return c.controller.UpdateMap(ctx, map2d)
		})
}

func (c *MapCommandSet) makeDeleteMapByIdCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"delete_map_by_id",
		cvalid.NewObjectSchema().
			WithRequiredProperty("map_id", cconv.String),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			return c.controller.DeleteMapById(ctx, args.GetAsString("map_id"))
		})
}
