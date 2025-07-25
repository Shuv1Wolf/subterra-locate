package test_persistence

import (
	"context"
	"testing"

	"github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/persistence"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

type BeaconsMemoryPersistenceTest struct {
	persistence *persistence.BeaconsMemoryPersistence
	fixture     *BeaconsPersistenceFixture
}

func newBeaconsMemoryPersistenceTest() *BeaconsMemoryPersistenceTest {
	persistence := persistence.NewBeaconsMemoryPersistence()
	persistence.Configure(context.Background(), cconf.NewEmptyConfigParams())

	fixture := NewBeaconsPersistenceFixture(persistence)

	return &BeaconsMemoryPersistenceTest{
		persistence: persistence,
		fixture:     fixture,
	}
}

func (c *BeaconsMemoryPersistenceTest) setup(t *testing.T) {
	err := c.persistence.Open(context.Background())
	if err != nil {
		t.Error("Failed to open persistence", err)
	}

	err = c.persistence.Clear(context.Background())
	if err != nil {
		t.Error("Failed to clear persistence", err)
	}
}

func (c *BeaconsMemoryPersistenceTest) teardown(t *testing.T) {
	err := c.persistence.Close(context.Background())
	if err != nil {
		t.Error("Failed to close persistence", err)
	}
}

func TestBeaconsMemoryPersistence(t *testing.T) {
	c := newBeaconsMemoryPersistenceTest()
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
