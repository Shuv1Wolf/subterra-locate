package test_clients1

import (
	"context"
	"testing"

	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/beacon-admin/clients/version1"
	controller1 "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/controllers/version1"
	persist "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/persistence"
	logic "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/service"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type beaconsCommandableHttpClientV1Test struct {
	persistence *persist.BeaconsMemoryPersistence
	service     *logic.BeaconsService
	controller  *controller1.BeaconCommandableGrpcControllerV1
	client      *clients1.BeaconsGrpcClientV1
	fixture     *BeaconsClientV1Fixture
	ctx         context.Context
}

func newBeaconsHttpClientV1Test() *beaconsCommandableHttpClientV1Test {
	ctx := context.Background()
	persistence := persist.NewBeaconsMemoryPersistence()
	persistence.Configure(ctx, cconf.NewEmptyConfigParams())

	service := logic.NewBeaconsService()
	service.Configure(ctx, cconf.NewEmptyConfigParams())

	httpConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.port", "3000",
		"connection.host", "localhost",
	)

	controller := controller1.NewBeaconCommandableGrpcControllerV1()
	controller.Configure(ctx, httpConfig)

	client := clients1.NewBeaconsGrpcClientV1()
	client.Configure(ctx, httpConfig)

	references := cref.NewReferencesFromTuples(ctx,
		cref.NewDescriptor("beacon-admin", "persistence", "memory", "default", "1.0"), persistence,
		cref.NewDescriptor("beacon-admin", "service", "default", "default", "1.0"), service,
		cref.NewDescriptor("beacon-admin", "controller", "grpc", "default", "1.0"), controller,
		cref.NewDescriptor("beacon-admin", "client", "grpc", "default", "1.0"), client,
	)
	service.SetReferences(ctx, references)
	controller.SetReferences(ctx, references)
	client.SetReferences(ctx, references)

	fixture := NewBeaconsClientV1Fixture(client)

	return &beaconsCommandableHttpClientV1Test{
		persistence: persistence,
		controller:  controller,
		service:     service,
		client:      client,
		fixture:     fixture,
		ctx:         ctx,
	}
}

func (c *beaconsCommandableHttpClientV1Test) setup(t *testing.T) {
	err := c.persistence.Open(c.ctx)
	if err != nil {
		t.Error("Failed to open persistence", err)
	}

	err = c.controller.Open(c.ctx)
	if err != nil {
		t.Error("Failed to open service", err)
	}

	err = c.client.Open(c.ctx)
	if err != nil {
		t.Error("Failed to open client", err)
	}

	err = c.persistence.Clear(c.ctx)
	if err != nil {
		t.Error("Failed to clear persistence", err)
	}
}

func (c *beaconsCommandableHttpClientV1Test) teardown(t *testing.T) {
	err := c.client.Close(c.ctx)
	if err != nil {
		t.Error("Failed to close client", err)
	}

	err = c.controller.Close(c.ctx)
	if err != nil {
		t.Error("Failed to close service", err)
	}

	err = c.persistence.Close(c.ctx)
	if err != nil {
		t.Error("Failed to close persistence", err)
	}
}

func TestBeaconsHttpClientV1(t *testing.T) {
	c := newBeaconsHttpClientV1Test()

	c.setup(t)

	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
