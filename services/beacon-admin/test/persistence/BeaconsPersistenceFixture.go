package test_persistence

import (
	"context"
	"testing"

	data "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/persistence"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

type BeaconsPersistenceFixture struct {
	BEACON1     *data.BeaconV1
	BEACON2     *data.BeaconV1
	BEACON3     *data.BeaconV1
	persistence persistence.IBeaconsPersistence
}

func NewBeaconsPersistenceFixture(persistence persistence.IBeaconsPersistence) *BeaconsPersistenceFixture {
	c := BeaconsPersistenceFixture{}

	c.BEACON1 = &data.BeaconV1{
		Id:    "1",
		Udi:   "00001",
		Type:  data.AltBeacon,
		Label: "TestBeacon1",
		X:     1.0,
		Y:     1.0,
		Z:     1.0,
		OrgId: "org$1001",
	}

	c.BEACON2 = &data.BeaconV1{
		Id:    "2",
		Udi:   "00002",
		Type:  data.IBeacon,
		OrgId: "org$1001",
		Label: "TestBeacon2",
		X:     1.0,
		Y:     1.0,
		Z:     1.0,
	}

	c.BEACON3 = &data.BeaconV1{
		Id:    "3",
		Udi:   "00003",
		Type:  data.AltBeacon,
		OrgId: "org$1001",
		Label: "TestBeacon3",
		X:     1.0,
		Y:     1.0,
		Z:     1.0,
	}

	c.persistence = persistence
	return &c
}

func (c *BeaconsPersistenceFixture) testCreateBeacons(t *testing.T) {
	// Create the first beacon
	beacon, err := c.persistence.Create(context.Background(), cdata.RequestContextV1{}, *c.BEACON1)
	assert.Nil(t, err)
	assert.NotEqual(t, data.BeaconV1{}, beacon)
	assert.Equal(t, c.BEACON1.Udi, beacon.Udi)
	assert.Equal(t, c.BEACON1.OrgId, beacon.OrgId)
	assert.Equal(t, c.BEACON1.Type, beacon.Type)
	assert.Equal(t, c.BEACON1.Label, beacon.Label)

	// Create the second beacon
	beacon, err = c.persistence.Create(context.Background(), cdata.RequestContextV1{}, *c.BEACON2)
	assert.Nil(t, err)
	assert.NotEqual(t, data.BeaconV1{}, beacon)
	assert.Equal(t, c.BEACON2.Udi, beacon.Udi)
	assert.Equal(t, c.BEACON2.OrgId, beacon.OrgId)
	assert.Equal(t, c.BEACON2.Type, beacon.Type)
	assert.Equal(t, c.BEACON2.Label, beacon.Label)

	// Create the third beacon
	beacon, err = c.persistence.Create(context.Background(), cdata.RequestContextV1{}, *c.BEACON3)
	assert.Nil(t, err)
	assert.NotEqual(t, data.BeaconV1{}, beacon)
	assert.Equal(t, c.BEACON3.Udi, beacon.Udi)
	assert.Equal(t, c.BEACON3.OrgId, beacon.OrgId)
	assert.Equal(t, c.BEACON3.Type, beacon.Type)
	assert.Equal(t, c.BEACON3.Label, beacon.Label)
}

func (c *BeaconsPersistenceFixture) TestCrudOperations(t *testing.T) {
	var beacon1 data.BeaconV1

	// Create items
	c.testCreateBeacons(t)

	// Get all beacons
	page, err := c.persistence.GetPageByFilter(context.Background(), cdata.RequestContextV1{},
		*cquery.NewEmptyFilterParams(), *cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 3)
	beacon1 = page.Data[0].Clone()

	// Update the beacon
	beacon1.Label = "ABC"
	beacon, err := c.persistence.Update(context.Background(), cdata.RequestContextV1{}, beacon1)
	assert.Nil(t, err)
	assert.NotEqual(t, data.BeaconV1{}, beacon)
	assert.Equal(t, beacon1.Id, beacon.Id)
	assert.Equal(t, "ABC", beacon.Label)

	// Get beacon by udi
	beacon, err = c.persistence.GetOneByUdi(context.Background(), cdata.RequestContextV1{}, beacon1.Udi)
	assert.Nil(t, err)
	assert.NotEqual(t, data.BeaconV1{}, beacon)
	assert.Equal(t, beacon1.Id, beacon.Id)

	// Delete the beacon
	beacon, err = c.persistence.DeleteById(context.Background(), cdata.RequestContextV1{}, beacon1.Id)
	assert.Nil(t, err)
	assert.NotEqual(t, data.BeaconV1{}, beacon)
	assert.Equal(t, beacon1.Id, beacon.Id)

	// Try to get deleted beacon
	beacon, err = c.persistence.GetOneById(context.Background(), cdata.RequestContextV1{}, beacon1.Id)
	assert.Nil(t, err)
	assert.Equal(t, data.BeaconV1{}, beacon)
}

func (c *BeaconsPersistenceFixture) TestGetWithFilters(t *testing.T) {
	// Create items
	c.testCreateBeacons(t)

	filter := *cquery.NewFilterParamsFromTuples(
		"id", "1",
	)
	// Filter by id
	page, err := c.persistence.GetPageByFilter(context.Background(), cdata.RequestContextV1{},
		filter,
		*cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 1)

	// Filter by udi
	filter = *cquery.NewFilterParamsFromTuples(
		"udi", "00002",
	)
	page, err = c.persistence.GetPageByFilter(
		context.Background(), cdata.RequestContextV1{},
		filter,
		*cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 1)

	// Filter by udis
	filter = *cquery.NewFilterParamsFromTuples(
		"udis", "00001,00003",
	)
	page, err = c.persistence.GetPageByFilter(
		context.Background(), cdata.RequestContextV1{},
		filter,
		*cquery.NewEmptyPagingParams())

	assert.Nil(t, err)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 2)

	// Filter by org_id
	filter = *cquery.NewFilterParamsFromTuples(
		"org_id", "org$1001",
	)
	page, err = c.persistence.GetPageByFilter(
		context.Background(), cdata.RequestContextV1{},
		filter,
		*cquery.NewEmptyPagingParams())

	assert.Nil(t, err)
	assert.Len(t, page.Data, 3)
}
