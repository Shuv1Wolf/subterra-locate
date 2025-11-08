package logic

import (
	"context"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	exec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type DeviceCommandSet struct {
	ccmd.CommandSet
	controller      IDeviceService
	deviceConvertor cconv.IJSONEngine[data.DeviceV1]
}

func NewDeviceCommandSet(controller IDeviceService) *DeviceCommandSet {
	c := &DeviceCommandSet{
		CommandSet:      *ccmd.NewCommandSet(),
		controller:      controller,
		deviceConvertor: cconv.NewDefaultCustomTypeJsonConvertor[data.DeviceV1](),
	}

	c.AddCommand(c.makeGetDevicesCommand())
	c.AddCommand(c.makeGetDeviceByIdCommand())
	c.AddCommand(c.makeCreateDeviceCommand())
	c.AddCommand(c.makeUpdateDeviceCommand())
	c.AddCommand(c.makeDeleteDeviceByIdCommand())

	return c
}

func (c *DeviceCommandSet) makeGetDevicesCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"get_devices",
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
			return c.controller.GetDevices(ctx, *reqctx, *filter, *paging)
		})
}

func (c *DeviceCommandSet) makeGetDeviceByIdCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"get_device_by_id",
		cvalid.NewObjectSchema().
			WithRequiredProperty("device_id", cconv.String).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.GetDeviceById(ctx, *reqctx, args.GetAsString("device_id"))
		})
}

func (c *DeviceCommandSet) makeCreateDeviceCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"create_device",
		cvalid.NewObjectSchema().
			WithRequiredProperty("device", data.NewDeviceV1Schema()).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {

			var device data.DeviceV1
			if _device, ok := args.GetAsObject("device"); ok {
				buf, err := cconv.JsonConverter.ToJson(_device)
				if err != nil {
					return nil, err
				}
				device, err = c.deviceConvertor.FromJson(buf)
				if err != nil {
					return nil, err
				}
			}
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.CreateDevice(ctx, *reqctx, device)
		})
}

func (c *DeviceCommandSet) makeUpdateDeviceCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"update_device",
		cvalid.NewObjectSchema().
			WithRequiredProperty("device", data.NewDeviceV1Schema()).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			var device data.DeviceV1
			if _device, ok := args.GetAsObject("device"); ok {
				buf, err := cconv.JsonConverter.ToJson(_device)
				if err != nil {
					return nil, err
				}
				device, err = c.deviceConvertor.FromJson(buf)
				if err != nil {
					return nil, err
				}
			}

			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.UpdateDevice(ctx, *reqctx, device)
		})
}

func (c *DeviceCommandSet) makeDeleteDeviceByIdCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"delete_device_by_id",
		cvalid.NewObjectSchema().
			WithRequiredProperty("device_id", cconv.String).
			WithOptionalProperty("reqctx", cdata.NewRequestContextV1Schema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			reqctx := cdata.NewRequestContextV1()
			if _val, ok := args.Get("reqctx"); ok {
				reqctx = cdata.NewRequestContextV1FromValue(_val)
			}
			return c.controller.DeleteDeviceById(ctx, *reqctx, args.GetAsString("device_id"))
		})
}
