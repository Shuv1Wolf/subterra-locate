package test_controllers

import (
	"context"
	"testing"
	"time"

	cclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"

	controllers "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/controllers/version1"
	data "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/persistence"
	logic "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/service"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	pipdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	tclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test"
	"github.com/stretchr/testify/assert"
)

type zonesGrpcControllerV1Test struct {
	ZONE1    *data.ZoneV1
	ZONE2    *data.ZoneV1
	persistence  *persistence.ZoneMemoryPersistence
	service         *logic.ZoneService
	controller      *controllers.ZoneCommandableGrpcControllerV1
	client          *tclients.TestCommandableGrpcClient
}

func newZonesGrpcControllerV1Test() *zonesGrpcControllerV1Test {
	ZONE1 := &data.ZoneV1{
		Id:        "1",
		MapId:     "map_1",
		OrgId:     "org_1",
		Name:      "TestZone1",
		PositionX: 10,
		PositionY: 10,
		Width:     100,
		Height:    100,
		Type:      "polygon",
		MaxDevice: 10,
		CreatedAt: time.Now(),
	}

	ZONE2 := &data.ZoneV1{
		Id:        "2",
		MapId:     "map_1",
		OrgId:     "org_1",
		Name:      "TestZone2",
		PositionX: 20,
		PositionY: 20,
		Width:     200,
		Height:    200,
		Type:      "circle",
		MaxDevice: 20,
		CreatedAt: time.Now(),
	}

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.port", "3001",
		"connection.host", "localhost",
	)

	persistence := persistence.NewZoneMemoryPersistence()
	persistence.Configure(context.Background(), cconf.NewEmptyConfigParams())

	service := logic.NewZoneService()
	service.Configure(context.Background(), cconf.NewEmptyConfigParams())

	controller := controllers.NewZoneCommandableGrpcControllerV1()
	controller.Configure(context.Background(), restConfig)

	client := tclients.NewTestCommandableGrpcClient("zone_processor.v1")
	client.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("zone-processor", "persistence", "memory", "default", "1.0"), persistence,
		cref.NewDescriptor("zone-processor", "service", "default", "default", "1.0"), service,
		cref.NewDescriptor("zone-processor", "controller", "http", "default", "1.0"), controller,
		cref.NewDescriptor("zone-processor", "client", "http", "default", "1.0"), client,
	)

	service.SetReferences(context.Background(), references)
	controller.SetReferences(context.Background(), references)

	return &zonesGrpcControllerV1Test{
		ZONE1:           ZONE1,
		ZONE2:           ZONE2,
		persistence:  persistence,
		controller:      controller,
		service:         service,
		client:          client,
	}
}

func (c *zonesGrpcControllerV1Test) setup(t *testing.T) {
	err := c.persistence.Open(context.Background())
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

	err = c.persistence.Clear(context.Background())
	if err != nil {
		t.Error("Failed to clear persistence", err)
	}
}

func (c *zonesGrpcControllerV1Test) teardown(t *testing.T) {
	err := c.client.Close(context.Background())
	if err != nil {
		t.Error("Failed to close client", err)
	}

	err = c.controller.Close(context.Background())
	if err != nil {
		t.Error("Failed to close service", err)
	}

	err = c.persistence.Close(context.Background())
	if err != nil {
		t.Error("Failed to close persistence", err)
	}
}

func (c *zonesGrpcControllerV1Test) testZoneCrudOperations(t *testing.T) {
	var zone1 data.ZoneV1

	// Create the first zone
	params := pipdata.NewAnyValueMapFromTuples(
		"zone", c.ZONE1,
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err := c.client.CallCommand(context.Background(), "create_zone", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	z, err := cclients.HandleHttpResponse[data.ZoneV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data.ZoneV1{}, z)
	assert.Equal(t, c.ZONE1.Name, z.Name)
	assert.Equal(t, c.ZONE1.OrgId, z.OrgId)

	// Create the second zone
	params = pipdata.NewAnyValueMapFromTuples(
		"zone", c.ZONE2,
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err = c.client.CallCommand(context.Background(), "create_zone", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	z, err = cclients.HandleHttpResponse[data.ZoneV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data.ZoneV1{}, z)
	assert.Equal(t, c.ZONE2.Name, z.Name)
	assert.Equal(t, c.ZONE2.OrgId, z.OrgId)

	// Get all zones
	params = pipdata.NewAnyValueMapFromTuples(
		"filter", cquery.NewEmptyFilterParams(),
		"paging", cquery.NewEmptyPagingParams(),
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err = c.client.CallCommand(context.Background(), "get_zones", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	page, err := cclients.HandleHttpResponse[cquery.DataPage[data.ZoneV1]](response)
	assert.Nil(t, err)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 2)
	zone1 = page.Data[0]

	// Update the zone
	zone1.Name = "ABC"
	params = pipdata.NewAnyValueMapFromTuples(
		"zone", zone1,
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err = c.client.CallCommand(context.Background(), "update_zone", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	z, err = cclients.HandleHttpResponse[data.ZoneV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data.ZoneV1{}, z)
	assert.Equal(t, zone1.Id, z.Id)
	assert.Equal(t, "ABC", z.Name)

	// Delete the zone
	params = pipdata.NewAnyValueMapFromTuples(
		"zone_id", zone1.Id,
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err = c.client.CallCommand(context.Background(), "delete_zone_by_id", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	z, err = cclients.HandleHttpResponse[data.ZoneV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data.ZoneV1{}, z)
	assert.Equal(t, zone1.Id, z.Id)

	// Try to get deleted zone
	params = pipdata.NewAnyValueMapFromTuples(
		"zone_id", zone1.Id,
		"reqctx", cdata.RequestContextV1{OrgId: "org_1"},
	)
	response, err = c.client.CallCommand(context.Background(), "get_zone_by_id", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	z, err = cclients.HandleHttpResponse[data.ZoneV1](response)
	assert.Nil(t, err)
	assert.Equal(t, data.ZoneV1{}, z)
}

func TestZonesCommmandableGrpcServiceV1(t *testing.T) {
	c := newZonesGrpcControllerV1Test()

	c.setup(t)
	t.Run("Zone CRUD Operations", c.testZoneCrudOperations)
	c.teardown(t)
}
