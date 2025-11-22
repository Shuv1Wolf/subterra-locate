package test_service

import (
	"context"
	"testing"
	"time"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/persistence"
	logic "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/service"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

type MapServiceTest struct {
	MAP1    *data.Map2dV1
	MAP2    *data.Map2dV1
	service *logic.MapService
}

func newMapServiceTest() *MapServiceTest {
	mapPersistence := persistence.NewMap2dMemoryPersistence()
	mapPersistence.Configure(context.Background(), cconf.NewEmptyConfigParams())

	srv := logic.NewMapService()
	srv.Configure(context.Background(), cconf.NewEmptyConfigParams())

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("geo-renderer", "persistence", "memory", "map-2d", "1.0"), mapPersistence,
	)
	srv.SetReferences(context.Background(), references)

	return &MapServiceTest{
		MAP1: &data.Map2dV1{
			Id:        "1",
			Name:      "TestMap1",
			SVG:       "<svg>map1</svg>",
			ScaleX:    1.0,
			ScaleY:    1.0,
			CreatedAt: time.Now(),
			OrgId:     "org_1",
			Width:     100,
			Height:    100,
			Level:     1,
		},
		MAP2: &data.Map2dV1{
			Id:        "2",
			Name:      "TestMap2",
			SVG:       "<svg>map2</svg>",
			ScaleX:    2.0,
			ScaleY:    2.0,
			CreatedAt: time.Now(),
			OrgId:     "org_1",
			Width:     200,
			Height:    200,
			Level:     2,
		},
		service: srv,
	}
}

func (c *MapServiceTest) TestMapCrudOperations(t *testing.T) {
	// Create one map
	m, err := c.service.CreateMap(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, *c.MAP1)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, c.MAP1.Id, m.Id)
	assert.Equal(t, c.MAP1.Name, m.Name)

	// Create another map
	m, err = c.service.CreateMap(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, *c.MAP2)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, c.MAP2.Id, m.Id)
	assert.Equal(t, c.MAP2.Name, m.Name)

	// Get all maps
	page, err := c.service.GetMaps(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, *cquery.NewEmptyFilterParams(), *cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 2)

	// Update the map
	m.Name = "Updated Map2"
	updatedMap, err := c.service.UpdateMap(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, m)
	assert.Nil(t, err)
	assert.NotNil(t, updatedMap)
	assert.Equal(t, m.Id, updatedMap.Id)
	assert.Equal(t, "Updated Map2", updatedMap.Name)

	// Delete map
	deletedMap, err := c.service.DeleteMapById(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, m.Id)
	assert.Nil(t, err)
	assert.NotNil(t, deletedMap)
	assert.Equal(t, m.Id, deletedMap.Id)

	// Try to get deleted map
	getMap, err := c.service.GetMapById(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, m.Id)
	assert.Nil(t, err)
	assert.Equal(t, data.Map2dV1{}, getMap)
}

func TestMapService(t *testing.T) {
	c := newMapServiceTest()
	t.Run("Map CRUD Operations", c.TestMapCrudOperations)
}
