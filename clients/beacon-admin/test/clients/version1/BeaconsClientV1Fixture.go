package test_clients1

import (
	"context"
	"testing"

	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/beacon-admin/clients/version1"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

type BeaconsClientV1Fixture struct {
	BEACON1 *data1.BeaconV1
	BEACON2 *data1.BeaconV1
	client  clients1.IBeaconsClientV1
	ctx     context.Context
}

func NewBeaconsClientV1Fixture(client clients1.IBeaconsClientV1) *BeaconsClientV1Fixture {
	c := &BeaconsClientV1Fixture{}

	c.BEACON1 = &data1.BeaconV1{
		Id:    "1",
		Udi:   "00001",
		Type:  data1.AltBeacon,
		OrgId: "1",
		Label: "TestBeacon1",
		X:     1.0,
		Y:     1.0,
		Z:     1.0,
	}

	c.BEACON2 = &data1.BeaconV1{
		Id:    "2",
		Udi:   "00002",
		Type:  data1.IBeacon,
		OrgId: "1",
		Label: "TestBeacon2",
		X:     1.0,
		Y:     1.0,
		Z:     1.0,
	}

	c.client = client
	c.ctx = context.Background()

	return c
}

func (c *BeaconsClientV1Fixture) testCreateBeacons(t *testing.T) {
	// Create the first beacon
	beacon, err := c.client.CreateBeacon(c.ctx, cdata.RequestContextV1{}, *c.BEACON1)
	assert.Nil(t, err)
	assert.NotNil(t, beacon)
	assert.Equal(t, c.BEACON1.Udi, beacon.Udi)
	assert.Equal(t, c.BEACON1.OrgId, beacon.OrgId)
	assert.Equal(t, c.BEACON1.Type, beacon.Type)
	assert.Equal(t, c.BEACON1.Label, beacon.Label)

	// Create the second beacon
	beacon, err = c.client.CreateBeacon(c.ctx, cdata.RequestContextV1{}, *c.BEACON2)
	assert.Nil(t, err)
	assert.NotNil(t, beacon)
	assert.Equal(t, c.BEACON2.Udi, beacon.Udi)
	assert.Equal(t, c.BEACON2.OrgId, beacon.OrgId)
	assert.Equal(t, c.BEACON2.Type, beacon.Type)
	assert.Equal(t, c.BEACON2.Label, beacon.Label)
}

func (c *BeaconsClientV1Fixture) TestCrudOperations(t *testing.T) {
	var beacon1 *data1.BeaconV1

	// Create items
	c.testCreateBeacons(t)

	// Get all beacons
	page, err := c.client.GetBeacons(c.ctx, cdata.RequestContextV1{}, cquery.NewEmptyFilterParams(), cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 2)
	beacon1 = &page.Data[0]

	// Update the beacon
	beacon1.Label = "ABC"
	beacon, err := c.client.UpdateBeacon(c.ctx, cdata.RequestContextV1{}, *beacon1)
	assert.Nil(t, err)
	assert.NotNil(t, beacon)
	assert.Equal(t, beacon1.Id, beacon.Id)
	assert.Equal(t, "ABC", beacon.Label)

	// Get beacon by udi
	beacon, err = c.client.GetBeaconByUdi(c.ctx, cdata.RequestContextV1{}, beacon1.Udi)
	assert.Nil(t, err)
	assert.NotNil(t, beacon)
	assert.Equal(t, beacon1.Id, beacon.Id)

	// Delete the beacon
	beacon, err = c.client.DeleteBeaconById(c.ctx, cdata.RequestContextV1{}, beacon1.Id)
	assert.Nil(t, err)
	assert.NotNil(t, beacon)
	assert.Equal(t, beacon1.Id, beacon.Id)

	// Try to get deleted beacon
	beacon, err = c.client.GetBeaconById(c.ctx, cdata.RequestContextV1{}, beacon1.Id)
	assert.Nil(t, err)
	assert.Empty(t, beacon)
}
