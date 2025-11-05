package clients1

import (
	"context"
	"reflect"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	mdata "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type GeoRendererMemoryClientV1 struct {
	maxPageSize int
	items       []data1.Map2dV1
	proto       reflect.Type
}

func NewGeoRendererMemoryClientV1(items []data1.Map2dV1) *GeoRendererMemoryClientV1 {
	c := &GeoRendererMemoryClientV1{
		maxPageSize: 100,
		items:       make([]data1.Map2dV1, 0),
		proto:       reflect.TypeOf(data1.Map2dV1{}),
	}
	c.items = append(c.items, items...)
	return c
}

func (c *GeoRendererMemoryClientV1) composeFilter(filter cquery.FilterParams) func(item data1.Map2dV1) bool {

	id := filter.GetAsString("id")
	orgId := filter.GetAsString("org_id")
	name := filter.GetAsString("name")

	return func(item data1.Map2dV1) bool {

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

func (c *GeoRendererMemoryClientV1) GetMaps(ctx context.Context,
	filter *cquery.FilterParams, paging *cquery.PagingParams) (page *cquery.DataPage[data1.Map2dV1], err error) {
	filterDevice := c.composeFilter(*filter)

	map2d := make([]data1.Map2dV1, 0)
	for _, v := range c.items {
		if filterDevice(v) {
			item := v
			map2d = append(map2d, item)
		}
	}

	skip := paging.GetSkip(-1)
	take := paging.GetTake((int64)(c.maxPageSize))
	var total int = 0

	if paging.Total {
		total = (len(map2d))
	}

	if skip > 0 {
		map2d = map2d[skip:]
	}

	if (int64)(len(map2d)) >= take {
		map2d = map2d[:take]
	}

	return cquery.NewDataPage(map2d, total), nil
}

func (c *GeoRendererMemoryClientV1) GetMapById(ctx context.Context,
	id string) (device *data1.Map2dV1, err error) {

	var item *data1.Map2dV1
	for _, v := range c.items {
		if v.Id == id {
			item = &v
			break
		}
	}

	return item, nil
}

func (c *GeoRendererMemoryClientV1) CreateMap(ctx context.Context,
	map2d data1.Map2dV1) (res *data1.Map2dV1, err error) {

	newItem := mdata.CloneObject(map2d, c.proto)
	item, _ := newItem.(data1.Map2dV1)
	mdata.GenerateObjectId(&newItem)

	c.items = append(c.items, item)
	return &item, nil
}

func (c *GeoRendererMemoryClientV1) UpdateMap(ctx context.Context,
	map2d data1.Map2dV1) (res *data1.Map2dV1, err error) {

	var index = -1
	for i, v := range c.items {
		if v.Id == map2d.Id {
			index = i
			break
		}
	}

	if index < 0 {
		return nil, nil
	}

	newItem := mdata.CloneObject(map2d, c.proto)
	item, _ := newItem.(data1.Map2dV1)
	c.items[index] = item
	return &item, nil
}

func (c *GeoRendererMemoryClientV1) DeleteMapById(ctx context.Context,
	id string) (res *data1.Map2dV1, err error) {

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
