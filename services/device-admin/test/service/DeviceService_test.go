package test_logic

import (
	"context"
	"testing"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	persist "github.com/Shuv1Wolf/subterra-locate/services/device-admin/persistence"
	logic "github.com/Shuv1Wolf/subterra-locate/services/device-admin/service"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

type DeviceServiceTest struct {
	DEVICE1     *data.DeviceV1
	DEVICE2     *data.DeviceV1
	persistence *persist.DeviceMemoryPersistence
	service     *logic.DeviceService
}

func newDeviceServiceTest() *DeviceServiceTest {
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

	persistence := persist.NewDeviceMemoryPersistence()
	persistence.Configure(context.Background(), cconf.NewEmptyConfigParams())

	service := logic.NewDeviceService()
	service.Configure(context.Background(), cconf.NewEmptyConfigParams())

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("device-admin", "persistence", "memory", "default", "1.0"), persistence,
		cref.NewDescriptor("device-admin", "service", "default", "default", "1.0"), service,
	)

	service.SetReferences(context.Background(), references)

	return &DeviceServiceTest{
		DEVICE1:     DEVICE1,
		DEVICE2:     DEVICE2,
		persistence: persistence,
		service:     service,
	}
}

func (c *DeviceServiceTest) setup(t *testing.T) {
	err := c.persistence.Open(context.Background())
	if err != nil {
		t.Error("Failed to open persistence", err)
	}

	err = c.persistence.Clear(context.Background())
	if err != nil {
		t.Error("Failed to clear persistence", err)
	}
}

func (c *DeviceServiceTest) teardown(t *testing.T) {
	err := c.persistence.Close(context.Background())
	if err != nil {
		t.Error("Failed to close persistence", err)
	}
}

func (c *DeviceServiceTest) testCrudOperations(t *testing.T) {
	var device1 data.DeviceV1

	// Create the first device
	device, err := c.service.CreateDevice(context.Background(), cdata.RequestContextV1{OrgId: "org_1001"}, c.DEVICE1.Clone())
	assert.Nil(t, err)
	assert.NotEqual(t, data.DeviceV1{}, device)
	assert.Equal(t, c.DEVICE1.Name, device.Name)
	assert.Equal(t, c.DEVICE1.OrgId, device.OrgId)
	assert.Equal(t, c.DEVICE1.Type, device.Type)
	assert.Equal(t, c.DEVICE1.MacAddress, device.MacAddress)

	// Create the second device
	device, err = c.service.CreateDevice(context.Background(), cdata.RequestContextV1{OrgId: "org_1001"}, c.DEVICE2.Clone())
	assert.Nil(t, err)
	assert.NotEqual(t, data.DeviceV1{}, device)
	assert.Equal(t, c.DEVICE2.Name, device.Name)
	assert.Equal(t, c.DEVICE2.OrgId, device.OrgId)
	assert.Equal(t, c.DEVICE2.Type, device.Type)
	assert.Equal(t, c.DEVICE2.MacAddress, device.MacAddress)

	// Get all beacons
	page, err := c.service.GetDevices(context.Background(), cdata.RequestContextV1{OrgId: "org_1001"}, *cquery.NewEmptyFilterParams(), *cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 2)
	device1 = page.Data[0].Clone()

	// Update the device
	device1.Name = "ABC"
	device, err = c.service.UpdateDevice(context.Background(), cdata.RequestContextV1{OrgId: "org_1001"}, device1)
	assert.Nil(t, err)
	assert.NotEqual(t, data.DeviceV1{}, device)
	assert.Equal(t, device1.Id, device.Id)
	assert.Equal(t, "ABC", device.Name)

	// Delete the device
	device, err = c.service.DeleteDeviceById(context.Background(), cdata.RequestContextV1{OrgId: "org_1001"}, device1.Id)
	assert.Nil(t, err)
	assert.NotEqual(t, data.DeviceV1{}, device)
	assert.Equal(t, device1.Id, device.Id)

	// Try to get deleted device
	device, err = c.service.GetDeviceById(context.Background(), cdata.RequestContextV1{OrgId: "org_1001"}, device1.Id)
	assert.Nil(t, err)
	assert.Equal(t, data.DeviceV1{}, device)
}

func TestDeviceService(t *testing.T) {
	c := newDeviceServiceTest()

	c.setup(t)
	t.Run("CRUD Operations", c.testCrudOperations)
	c.teardown(t)
}
