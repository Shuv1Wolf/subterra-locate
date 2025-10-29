package test_controllers

import (
	"context"
	"testing"

	data "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/device-admin/persistence"
	logic "github.com/Shuv1Wolf/subterra-locate/services/device-admin/service"
	controllers "github.com/Shuv1Wolf/subterra-locate/services/device-admin/controllers/version1"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/stretchr/testify/assert"
)

type DeviceGrpcControllerV1Test struct {
	DEVICE1     *data.DeviceV1
	DEVICE2     *data.DeviceV1
	persistence *persistence.DeviceMemoryPersistence
	service     *logic.DeviceService
	controller  *controllers.DeviceCommandableGrpcControllerV1
}

func newDeviceGrpcControllerV1Test() *DeviceGrpcControllerV1Test {
	DEVICE1 := &data.DeviceV1{
		Id:         "1",
		Name:       "TestDevice1",
		Type:       "unknown",
		Model:      "test_model_1",
		OrgId:      "org_1001",
		Enabled:    true,
		MacAddress: "00:00:00:00:00:01",
		IpAddress:  "127.0.0.1",
	}

	DEVICE2 := &data.DeviceV1{
		Id:         "2",
		Name:       "TestDevice2",
		Type:       "smartphone",
		Model:      "test_model_2",
		OrgId:      "org_1001",
		Enabled:    true,
		MacAddress: "00:00:00:00:00:02",
		IpAddress:  "127.0.0.2",
	}

	persistence := persistence.NewDeviceMemoryPersistence()
	service := logic.NewDeviceService()
	controller := controllers.NewDeviceCommandableGrpcControllerV1()

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("device-admin", "persistence", "memory", "default", "1.0"), persistence,
		cref.NewDescriptor("device-admin", "service", "default", "default", "1.0"), service,
		cref.NewDescriptor("device-admin", "controller", "grpc", "default", "1.0"), controller,
	)

	service.SetReferences(context.Background(), references)
	controller.SetReferences(context.Background(), references)

	return &DeviceGrpcControllerV1Test{
		DEVICE1:     DEVICE1,
		DEVICE2:     DEVICE2,
		persistence: persistence,
		service:     service,
		controller:  controller,
	}
}

func (c *DeviceGrpcControllerV1Test) setup(t *testing.T) {
	err := c.persistence.Open(context.Background())
	if err != nil {
		t.Error("Failed to open persistence", err)
	}
}

func (c *DeviceGrpcControllerV1Test) teardown(t *testing.T) {
	err := c.persistence.Close(context.Background())
	if err != nil {
		t.Error("Failed to close persistence", err)
	}
}

func (c *DeviceGrpcControllerV1Test) testCrudOperations(t *testing.T) {
	// Create one device
	params := cexec.NewParametersFromTuples("device", c.DEVICE1)
	res, err := c.service.GetCommandSet().Execute(context.Background(), "create_device", params)
	assert.Nil(t, err)
	device := res.(data.DeviceV1)
	assert.NotNil(t, device)
	assert.Equal(t, c.DEVICE1.Name, device.Name)
	assert.Equal(t, c.DEVICE1.OrgId, device.OrgId)

	// Create another device
	params = cexec.NewParametersFromTuples("device", c.DEVICE2)
	res, err = c.service.GetCommandSet().Execute(context.Background(), "create_device", params)
	assert.Nil(t, err)
	device = res.(data.DeviceV1)
	assert.NotNil(t, device)
	assert.Equal(t, c.DEVICE2.Name, device.Name)
	assert.Equal(t, c.DEVICE2.OrgId, device.OrgId)

	// Get all devices
	res, err = c.service.GetCommandSet().Execute(context.Background(), "get_devices", cexec.NewEmptyParameters())
	assert.Nil(t, err)
	page := res.(map[string]any)
	assert.NotNil(t, page)
	assert.Len(t, page["data"], 2)

	// Update the device
	device.Name = "Updated Name"
	params = cexec.NewParametersFromTuples("device", device)
	res, err = c.service.GetCommandSet().Execute(context.Background(), "update_device", params)
	assert.Nil(t, err)
	updatedDevice := res.(data.DeviceV1)
	assert.NotNil(t, updatedDevice)
	assert.Equal(t, "Updated Name", updatedDevice.Name)

	// Delete the device
	params = cexec.NewParametersFromTuples("device_id", device.Id)
	res, err = c.service.GetCommandSet().Execute(context.Background(), "delete_device_by_id", params)
	assert.Nil(t, err)
	deletedDevice := res.(data.DeviceV1)
	assert.NotNil(t, deletedDevice)
	assert.Equal(t, device.Id, deletedDevice.Id)

	// Try to get deleted device
	params = cexec.NewParametersFromTuples("device_id", device.Id)
	res, err = c.service.GetCommandSet().Execute(context.Background(), "get_device_by_id", params)
	assert.Nil(t, err)
	assert.Equal(t, data.DeviceV1{}, res)
}

func TestDeviceGrpcControllerV1(t *testing.T) {
	c := newDeviceGrpcControllerV1Test()
	c.setup(t)
	t.Run("CRUD Operations", c.testCrudOperations)
	c.teardown(t)
}
