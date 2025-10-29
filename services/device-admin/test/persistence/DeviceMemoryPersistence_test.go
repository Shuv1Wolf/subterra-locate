package test_persistence

import (
	"context"
	"testing"

	"github.com/Shuv1Wolf/subterra-locate/services/device-admin/persistence"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

type DeviceMemoryPersistenceTest struct {
	persistence *persistence.DeviceMemoryPersistence
	fixture     *DevicePersistenceFixture
}

func newDeviceMemoryPersistenceTest() *DeviceMemoryPersistenceTest {
	persistence := persistence.NewDeviceMemoryPersistence()
	persistence.Configure(context.Background(), cconf.NewEmptyConfigParams())

	fixture := NewDevicePersistenceFixture(persistence)

	return &DeviceMemoryPersistenceTest{
		persistence: persistence,
		fixture:     fixture,
	}
}

func (c *DeviceMemoryPersistenceTest) setup(t *testing.T) {
	err := c.persistence.Open(context.Background())
	if err != nil {
		t.Error("Failed to open persistence", err)
	}

	err = c.persistence.Clear(context.Background())
	if err != nil {
		t.Error("Failed to clear persistence", err)
	}
}

func (c *DeviceMemoryPersistenceTest) teardown(t *testing.T) {
	err := c.persistence.Close(context.Background())
	if err != nil {
		t.Error("Failed to close persistence", err)
	}
}

func TestDeviceMemoryPersistence(t *testing.T) {
	c := newDeviceMemoryPersistenceTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)

	c.setup(t)
	t.Run("Get With Filters", c.fixture.TestGetWithFilters)
	c.teardown(t)
}
