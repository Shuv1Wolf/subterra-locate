package logic

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/device-admin/persistence"
	"github.com/Shuv1Wolf/subterra-locate/services/device-admin/publisher"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type DeviceService struct {
	persistence  persistence.IDevicePersistence
	commandSet   *DeviceCommandSet
	deviceEvents publisher.IPublisher
}

func NewDeviceService() *DeviceService {
	c := &DeviceService{}
	return c
}

func (c *DeviceService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	// Read configuration parameters here...
}

func (c *DeviceService) GetCommandSet() *ccmd.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewDeviceCommandSet(c)
	}
	return &c.commandSet.CommandSet
}

func (c *DeviceService) SetReferences(ctx context.Context, references cref.IReferences) {
	res, err := references.GetOneRequired(
		cref.NewDescriptor("device-admin", "persistence", "*", "*", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.persistence = res.(persistence.IDevicePersistence)

	res = references.GetOneOptional(
		cref.NewDescriptor("device-admin", "publisher", "nats", "device-events", "1.0"),
	)
	c.deviceEvents = res.(publisher.IPublisher)
}

func (c *DeviceService) GetDevices(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.DeviceV1], error) {
	return c.persistence.GetPageByFilter(ctx, filter, paging)
}

func (c *DeviceService) GetDeviceById(ctx context.Context,
	deviceId string) (data.DeviceV1, error) {

	return c.persistence.GetOneById(ctx, deviceId)
}

func (c *DeviceService) CreateDevice(ctx context.Context,
	device data.DeviceV1) (data.DeviceV1, error) {

	if device.Type == "" {
		device.Type = data.Unknown
	}

	b, err := c.persistence.Create(ctx, device)
	if err != nil {
		return b, err
	}

	if c.deviceEvents != nil {
		err = c.deviceEvents.SendDeviceCreatedEvent(ctx, b.Id)
		if err != nil {
			return b, err
		}
	}

	return b, nil
}

func (c *DeviceService) UpdateDevice(ctx context.Context,
	device data.DeviceV1) (data.DeviceV1, error) {

	if device.Type == "" {
		device.Type = data.Unknown
	}

	b, err := c.persistence.Update(ctx, device)
	if err != nil {
		return b, err
	}

	if c.deviceEvents != nil {
		err = c.deviceEvents.SendDeviceChangedEvent(ctx, b.Id)
		if err != nil {
			return b, err
		}
	}

	return b, err
}

func (c *DeviceService) DeleteDeviceById(ctx context.Context,
	devcieId string) (data.DeviceV1, error) {

	b, err := c.persistence.DeleteById(ctx, devcieId)
	if err != nil {
		return b, err
	}

	if c.deviceEvents != nil {
		err = c.deviceEvents.SendDeviceDeletedEvent(ctx, b.Id)
		if err != nil {
			return b, err
		}
	}

	return b, err
}
