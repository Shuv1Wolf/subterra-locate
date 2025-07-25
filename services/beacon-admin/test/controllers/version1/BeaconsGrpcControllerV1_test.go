package test_services1

import (
	"context"
	"testing"

	cclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"

	controllers1 "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/controllers/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	persist "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/persistence"
	logic "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/service"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	tclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test"
	"github.com/stretchr/testify/assert"
)

type beaconsGrpcControllerV1Test struct {
	BEACON1     *data1.BeaconV1
	BEACON2     *data1.BeaconV1
	persistence *persist.BeaconsMemoryPersistence
	service     *logic.BeaconsService
	controller  *controllers1.BeaconCommandableGrpcControllerV1
	client      *tclients.TestCommandableGrpcClient
}

func newbeaconsGrpcControllerV1Test() *beaconsGrpcControllerV1Test {
	BEACON1 := &data1.BeaconV1{
		Id:     "1",
		Udi:    "00001",
		Type:   data1.AltBeacon,
		SiteId: "1",
		Label:  "TestBeacon1",
		X:      1.0,
		Y:      1.0,
		Z:      1.0,
	}

	BEACON2 := &data1.BeaconV1{
		Id:     "2",
		Udi:    "00002",
		Type:   data1.IBeacon,
		SiteId: "1",
		Label:  "TestBeacon2",
		X:      1.0,
		Y:      1.0,
		Z:      1.0,
	}

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.port", "3000",
		"connection.host", "localhost",
	)

	persistence := persist.NewBeaconsMemoryPersistence()
	persistence.Configure(context.Background(), cconf.NewEmptyConfigParams())

	service := logic.NewBeaconsService()
	service.Configure(context.Background(), cconf.NewEmptyConfigParams())

	controller := controllers1.NewBeaconCommandableGrpcControllerV1()
	controller.Configure(context.Background(), restConfig)

	client := tclients.NewTestCommandableGrpcClient("beacon.admin.v1")
	client.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("beacon-admin", "persistence", "memory", "default", "1.0"), persistence,
		cref.NewDescriptor("beacon-admin", "service", "default", "default", "1.0"), service,
		cref.NewDescriptor("beacon-admin", "controller", "http", "default", "1.0"), controller,
		cref.NewDescriptor("beacon-admin", "client", "http", "default", "1.0"), client,
	)

	service.SetReferences(context.Background(), references)
	controller.SetReferences(context.Background(), references)

	return &beaconsGrpcControllerV1Test{
		BEACON1:     BEACON1,
		BEACON2:     BEACON2,
		persistence: persistence,
		controller:  controller,
		service:     service,
		client:      client,
	}
}

func (c *beaconsGrpcControllerV1Test) setup(t *testing.T) {
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

func (c *beaconsGrpcControllerV1Test) teardown(t *testing.T) {
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

func (c *beaconsGrpcControllerV1Test) testCrudOperations(t *testing.T) {
	var beacon1 data1.BeaconV1

	// Create the first beacon
	params := cdata.NewAnyValueMapFromTuples(
		"beacon", c.BEACON1.Clone(),
	)
	response, err := c.client.CallCommand(context.Background(), "create_beacon", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	beacon, err := cclients.HandleHttpResponse[data1.BeaconV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data1.BeaconV1{}, beacon)
	assert.Equal(t, c.BEACON1.Udi, beacon.Udi)
	assert.Equal(t, c.BEACON1.SiteId, beacon.SiteId)
	assert.Equal(t, c.BEACON1.Type, beacon.Type)
	assert.Equal(t, c.BEACON1.Label, beacon.Label)

	// Create the second beacon
	params = cdata.NewAnyValueMapFromTuples(
		"beacon", c.BEACON2.Clone(),
	)
	response, err = c.client.CallCommand(context.Background(), "create_beacon", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	beacon, err = cclients.HandleHttpResponse[data1.BeaconV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data1.BeaconV1{}, beacon)
	assert.Equal(t, c.BEACON2.Udi, beacon.Udi)
	assert.Equal(t, c.BEACON2.SiteId, beacon.SiteId)
	assert.Equal(t, c.BEACON2.Type, beacon.Type)
	assert.Equal(t, c.BEACON2.Label, beacon.Label)

	// Get all beacons
	params = cdata.NewAnyValueMapFromTuples(
		"filter", cquery.NewEmptyFilterParams(),
		"paging", cquery.NewEmptyFilterParams(),
	)
	response, err = c.client.CallCommand(context.Background(), "get_beacons", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	page, err := cclients.HandleHttpResponse[cquery.DataPage[data1.BeaconV1]](response)
	assert.Nil(t, err)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 2)
	beacon1 = page.Data[0].Clone()

	// Update the beacon
	beacon1.Label = "ABC"
	params = cdata.NewAnyValueMapFromTuples(
		"beacon", beacon1,
	)
	response, err = c.client.CallCommand(context.Background(), "update_beacon", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	beacon, err = cclients.HandleHttpResponse[data1.BeaconV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data1.BeaconV1{}, beacon)
	assert.Equal(t, c.BEACON1.Id, beacon.Id)
	assert.Equal(t, "ABC", beacon.Label)

	// Get beacon by udi
	params = cdata.NewAnyValueMapFromTuples(
		"udi", beacon1.Udi,
	)
	response, err = c.client.CallCommand(context.Background(), "get_beacon_by_udi", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	beacon, err = cclients.HandleHttpResponse[data1.BeaconV1](response)
	assert.Nil(t, err)
	assert.NotEqual(t, data1.BeaconV1{}, beacon)
	assert.Equal(t, c.BEACON1.Id, beacon.Id)

	// Delete the beacon
	params = cdata.NewAnyValueMapFromTuples(
		"beacon_id", beacon1.Id,
	)
	response, err = c.client.CallCommand(context.Background(), "delete_beacon_by_id", params)
	assert.Nil(t, err)

	beacon, err = cclients.HandleHttpResponse[data1.BeaconV1](response)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	assert.NotEqual(t, data1.BeaconV1{}, beacon)
	assert.Equal(t, c.BEACON1.Id, beacon.Id)

	// Try to get deleted beacon
	params = cdata.NewAnyValueMapFromTuples(
		"beacon_id", beacon1.Id,
	)
	response, err = c.client.CallCommand(context.Background(), "get_beacon_by_id", params)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	beacon, err = cclients.HandleHttpResponse[data1.BeaconV1](response)
	assert.Nil(t, err)
	assert.Equal(t, data1.BeaconV1{}, beacon)
}

func TestBeaconsCommmandableGrpcpServiceV1(t *testing.T) {
	c := newbeaconsGrpcControllerV1Test()

	c.setup(t)
	t.Run("CRUD Operations", c.testCrudOperations)
	c.teardown(t)
}
