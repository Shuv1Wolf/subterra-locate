package test_clients1

import (
	"testing"

	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/beacon-admin/clients/version1"
)

type beaconsMemoryClientV1Test struct {
	client  *clients1.BeaconsMemoryClientV1
	fixture *BeaconsClientV1Fixture
}

func newBeaconsMemoryClientV1Test() *beaconsMemoryClientV1Test {

	return &beaconsMemoryClientV1Test{}
}

func (c *beaconsMemoryClientV1Test) setup(t *testing.T) {
	c.client = clients1.NewBeaconsMemoryClientV1(nil)
	c.fixture = NewBeaconsClientV1Fixture(c.client)
}

func (c *beaconsMemoryClientV1Test) teardown(t *testing.T) {
	c.client = nil
	c.fixture = nil
}

func TestBeaconsMemoryClientV1(t *testing.T) {
	c := newBeaconsMemoryClientV1Test()

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
