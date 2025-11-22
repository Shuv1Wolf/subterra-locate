package test_controllers

import (
	"context"
	"testing"
	"time"

	cclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	controllers "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/controllers/version1"
	data "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/persistence"
	logic "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/service"
	pipdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	tclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test"
	"github.com/stretchr/testify/assert"
)

type mapsGrpcControllerV1Test struct {
	MAP1           *data.Map2dV1
	MAP2           *data.Map2dV1
	mapPersistence *persistence.Map2dMemoryPersistence
	service        *logic.MapService
	controller     *controllers.MapCommandableGrpcControllerV1
	client         *tclients.TestCommandableGrpcClient
}

func newMapsGrpcControllerV1Test() *mapsGrpcControllerV1Test {
	MAP1 := &data.Map2dV1{
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
	}

	MAP2 := &data.Map2dV1{
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
	}

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.port", "3000",
		"connection.host", "localhost",
	)

	mapPersistence := persistence.NewMap2dMemoryPersistence()
	mapPersistence.Configure(context.Background(), cconf.NewEmptyConfigParams())

	service := logic.NewMapService()
	service.Configure(context.Background(), cconf.NewEmptyConfigParams())

	controller := controllers.NewMapCommandableGrpcControllerV1()
	controller.Configure(context.Background(), restConfig)

	client := tclients.NewTestCommandableGrpcClient("geo.renderer.v1")
	client.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("geo-renderer", "persistence", "memory", "map-2d", "1.0"), mapPersistence,
		cref.NewDescriptor("geo-renderer", "service", "default", "default", "1.0"), service,
		cref.NewDescriptor("geo-renderer", "controller", "http", "default", "1.0"), controller,
		cref.NewDescriptor("geo-renderer", "client", "http", "default", "1.0"), client,
	)

	service.SetReferences(context.Background(), references)
	controller.SetReferences(context.Background(), references)

	return &mapsGrpcControllerV1Test{
		MAP1:           MAP1,
		MAP2:           MAP2,
		mapPersistence: mapPersistence,
		controller:     controller,
		service:        service,
		client:         client,
	}
}

func (c *mapsGrpcControllerV1Test) setup(t *testing.T) {
	err := c.mapPersistence.Open(context.Background())
	if err != nil {
		t.Error("Failed to open persistence", err)
	}

	err = c.controller.Open(context.Background())
	if err != nil {
		t.Error("Failed to open service", err)
	}

	err = c.client.Open(context.Background())
	if err != nil {
		t.Error("Failed to open client", err)
	}

	err = c.mapPersistence.Clear(context.Background())
	if err != nil {
		t.Error("Failed to clear persistence", err)
	}
}

func (c *mapsGrpcControllerV1Test) teardown(t *testing.T) {
	err := c.client.Close(context.Background())
	if err != nil {
		t.Error("Failed to close client", err)
	}

	err = c.controller.Close(context.Background())
	if err != nil {
		t.Error("Failed to close service", err)
	}

	err = c.mapPersistence.Close(context.Background())
	if err != nil {
		t.Error("Failed to close persistence", err)
	}
}

func (c *mapsGrpcControllerV1Test) testMapCrudOperations(t *testing.T) {
	var map1 data.Map2dV1

	// Create the first map
	params := pipdata.NewAnyValueMapFromTuples(
		"map", c.MAP1,
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err := c.client.CallCommand(context.Background(), "create_map", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	m, err := cclients.HandleHttpResponse[data.Map2dV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data.Map2dV1{}, m)
	assert.Equal(t, c.MAP1.Name, m.Name)
	assert.Equal(t, c.MAP1.OrgId, m.OrgId)

	// Create the second map
	params = pipdata.NewAnyValueMapFromTuples(
		"map", c.MAP2,
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err = c.client.CallCommand(context.Background(), "create_map", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	m, err = cclients.HandleHttpResponse[data.Map2dV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data.Map2dV1{}, m)
	assert.Equal(t, c.MAP2.Name, m.Name)
	assert.Equal(t, c.MAP2.OrgId, m.OrgId)

	// Get all maps
	params = pipdata.NewAnyValueMapFromTuples(
		"filter", cquery.NewEmptyFilterParams(),
		"paging", cquery.NewEmptyPagingParams(),
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err = c.client.CallCommand(context.Background(), "get_maps", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	page, err := cclients.HandleHttpResponse[cquery.DataPage[data.Map2dV1]](response)
	assert.Nil(t, err)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 2)
	map1 = page.Data[0]

	// Update the map
	map1.Name = "ABC"
	params = pipdata.NewAnyValueMapFromTuples(
		"map", map1,
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err = c.client.CallCommand(context.Background(), "update_map", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	m, err = cclients.HandleHttpResponse[data.Map2dV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data.Map2dV1{}, m)
	assert.Equal(t, map1.Id, m.Id)
	assert.Equal(t, "ABC", m.Name)

	// Delete the map
	params = pipdata.NewAnyValueMapFromTuples(
		"map_id", map1.Id,
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err = c.client.CallCommand(context.Background(), "delete_map_by_id", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	m, err = cclients.HandleHttpResponse[data.Map2dV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data.Map2dV1{}, m)
	assert.Equal(t, map1.Id, m.Id)

	// Try to get deleted map
	params = pipdata.NewAnyValueMapFromTuples(
		"map_id", map1.Id,
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err = c.client.CallCommand(context.Background(), "get_map_by_id", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	m, err = cclients.HandleHttpResponse[data.Map2dV1](response)
	assert.Nil(t, err)
	assert.Equal(t, data.Map2dV1{}, m)
}

func TestMapsCommmandableGrpcServiceV1(t *testing.T) {
	c := newMapsGrpcControllerV1Test()

	c.setup(t)
	t.Run("Map CRUD Operations", c.testMapCrudOperations)
	c.teardown(t)
}
