package test_persistence

import (
	"context"
	"testing"

	data "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/device-admin/persistence"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

type DevicePersistenceFixture struct {
	DEVICE1     *data.DeviceV1
	DEVICE2     *data.DeviceV1
	DEVICE3     *data.DeviceV1
	persistence persistence.IDevicePersistence
}

func NewDevicePersistenceFixture(persistence persistence.IDevicePersistence) *DevicePersistenceFixture {
	c := DevicePersistenceFixture{}

	c.DEVICE1 = &data.DeviceV1{
		Id:         "1",
		Name:       "TestDevice1",
		Type:       "unknown",
		Model:      "test_model_1",
		OrgId:      "org_1001",
		Enabled:    true,
		MacAddress: "00:00:00:00:00:01",
		IpAddress:  "127.0.0.1",
	}

	c.DEVICE2 = &data.DeviceV1{
		Id:         "2",
		Name:       "TestDevice2",
		Type:       "smartphone",
		Model:      "test_model_2",
		OrgId:      "org_1001",
		Enabled:    true,
		MacAddress: "00:00:00:00:00:02",
		IpAddress:  "127.0.0.2",
	}

	c.DEVICE3 = &data.DeviceV1{
		Id:         "3",
		Name:       "TestDevice3",
		Type:       "unknown",
		Model:      "test_model_3",
		OrgId:      "org_1002",
		Enabled:    false,
		MacAddress: "00:00:00:00:00:03",
		IpAddress:  "127.0.0.3",
	}

	c.persistence = persistence
	return &c
}

func (c *DevicePersistenceFixture) testCreateDevices(t *testing.T) {
	// Create the first device
	device, err := c.persistence.Create(context.Background(), *c.DEVICE1)
	assert.Nil(t, err)
	assert.NotEqual(t, data.DeviceV1{}, device)
	assert.Equal(t, c.DEVICE1.Name, device.Name)
	assert.Equal(t, c.DEVICE1.OrgId, device.OrgId)
	assert.Equal(t, c.DEVICE1.Type, device.Type)
	assert.Equal(t, c.DEVICE1.MacAddress, device.MacAddress)

	// Create the second device
	device, err = c.persistence.Create(context.Background(), *c.DEVICE2)
	assert.Nil(t, err)
	assert.NotEqual(t, data.DeviceV1{}, device)
	assert.Equal(t, c.DEVICE2.Name, device.Name)
	assert.Equal(t, c.DEVICE2.OrgId, device.OrgId)
	assert.Equal(t, c.DEVICE2.Type, device.Type)
	assert.Equal(t, c.DEVICE2.MacAddress, device.MacAddress)

	// Create the third device
	device, err = c.persistence.Create(context.Background(), *c.DEVICE3)
	assert.Nil(t, err)
	assert.NotEqual(t, data.DeviceV1{}, device)
	assert.Equal(t, c.DEVICE3.Name, device.Name)
	assert.Equal(t, c.DEVICE3.OrgId, device.OrgId)
	assert.Equal(t, c.DEVICE3.Type, device.Type)
	assert.Equal(t, c.DEVICE3.MacAddress, device.MacAddress)
}

func (c *DevicePersistenceFixture) TestCrudOperations(t *testing.T) {
	var device1 data.DeviceV1

	// Create items
	c.testCreateDevices(t)

	// Get all beacons
	page, err := c.persistence.GetPageByFilter(context.Background(),
		*cquery.NewEmptyFilterParams(), *cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 3)
	device1 = page.Data[0].Clone()

	// Update the device
	device1.Name = "ABC"
	device, err := c.persistence.Update(context.Background(), device1)
	assert.Nil(t, err)
	assert.NotEqual(t, data.DeviceV1{}, device)
	assert.Equal(t, device1.Id, device.Id)
	assert.Equal(t, "ABC", device.Name)

	// Delete the device
	device, err = c.persistence.DeleteById(context.Background(), device1.Id)
	assert.Nil(t, err)
	assert.NotEqual(t, data.DeviceV1{}, device)
	assert.Equal(t, device1.Id, device.Id)

	// Try to get deleted device
	device, err = c.persistence.GetOneById(context.Background(), device1.Id)
	assert.Nil(t, err)
	assert.Equal(t, data.DeviceV1{}, device)
}

func (c *DevicePersistenceFixture) TestGetWithFilters(t *testing.T) {
	// Create items
	c.testCreateDevices(t)

	filter := *cquery.NewFilterParamsFromTuples(
		"id", "1",
	)
	// Filter by id
	page, err := c.persistence.GetPageByFilter(context.Background(),
		filter,
		*cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 1)

	// Filter by mac
	filter = *cquery.NewFilterParamsFromTuples(
		"mac_address", "00:00:00:00:00:02",
	)
	page, err = c.persistence.GetPageByFilter(
		context.Background(),
		filter,
		*cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 1)

	// Filter by org_id
	filter = *cquery.NewFilterParamsFromTuples(
		"org_id", "org_1001",
	)
	page, err = c.persistence.GetPageByFilter(
		context.Background(),
		filter,
		*cquery.NewEmptyPagingParams())

	assert.Nil(t, err)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 2)
}
