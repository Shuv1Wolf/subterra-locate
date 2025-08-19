package clients1

import (
	"context"
	"reflect"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	mdata "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type DeviceMemoryClientV1 struct {
	maxPageSize int
	items       []data1.DeviceV1
	proto       reflect.Type
}

func NewDeviceMemoryClientV1(items []data1.DeviceV1) *DeviceMemoryClientV1 {
	c := &DeviceMemoryClientV1{
		maxPageSize: 100,
		items:       make([]data1.DeviceV1, 0),
		proto:       reflect.TypeOf(data1.DeviceV1{}),
	}
	c.items = append(c.items, items...)
	return c
}

func (c *DeviceMemoryClientV1) composeFilter(filter cquery.FilterParams) func(item data1.DeviceV1) bool {

	id := filter.GetAsString("id")
	orgId := filter.GetAsString("org_id")
	name := filter.GetAsString("name")

	return func(item data1.DeviceV1) bool {

		if id != "" && item.Id != id {
			return false
		}
		if orgId != "" && item.OrgId != orgId {
			return false
		}
		if name != "" && item.Name != name {
			return false
		}
		return true
	}
}

func (c *DeviceMemoryClientV1) GetDevices(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page *cquery.DataPage[data1.DeviceV1], err error) {
	filterDevice := c.composeFilter(filter)

	Device := make([]data1.DeviceV1, 0)
	for _, v := range c.items {
		if filterDevice(v) {
			item := v
			Device = append(Device, item)
		}
	}

	skip := paging.GetSkip(-1)
	take := paging.GetTake((int64)(c.maxPageSize))
	var total int = 0

	if paging.Total {
		total = (len(Device))
	}

	if skip > 0 {
		Device = Device[skip:]
	}

	if (int64)(len(Device)) >= take {
		Device = Device[:take]
	}

	return cquery.NewDataPage(Device, total), nil
}

func (c *DeviceMemoryClientV1) GetDEviceById(ctx context.Context,
	id string) (device *data1.DeviceV1, err error) {

	var item *data1.DeviceV1
	for _, v := range c.items {
		if v.Id == id {
			item = &v
			break
		}
	}

	return item, nil
}

func (c *DeviceMemoryClientV1) CreateDevice(ctx context.Context,
	device data1.DeviceV1) (res *data1.DeviceV1, err error) {

	newItem := mdata.CloneObject(device, c.proto)
	item, _ := newItem.(data1.DeviceV1)
	mdata.GenerateObjectId(&newItem)

	c.items = append(c.items, item)
	return &item, nil
}

func (c *DeviceMemoryClientV1) UpdateDevice(ctx context.Context,
	device data1.DeviceV1) (res *data1.DeviceV1, err error) {

	var index = -1
	for i, v := range c.items {
		if v.Id == device.Id {
			index = i
			break
		}
	}

	if index < 0 {
		return nil, nil
	}

	newItem := mdata.CloneObject(device, c.proto)
	item, _ := newItem.(data1.DeviceV1)
	c.items[index] = item
	return &item, nil
}

func (c *DeviceMemoryClientV1) DeleteDeviceById(ctx context.Context,
	id string) (res *data1.DeviceV1, err error) {

	var index = -1
	for i, v := range c.items {
		if v.Id == id {
			index = i
			break
		}
	}

	if index < 0 {
		return nil, nil
	}

	var item = c.items[index]

	if index == len(c.items) {
		c.items = c.items[:index-1]
	} else {
		c.items = append(c.items[:index], c.items[index+1:]...)
	}
	return &item, nil
}
