package test_persistence

import (
	"context"
	"testing"
	"time"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/persistence"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

type ZonePersistenceFixture struct {
	ZONE1       *data.ZoneV1
	ZONE2       *data.ZoneV1
	ZONE3       *data.ZoneV1
	persistence persistence.IZonePersistence
}

func NewZonePersistenceFixture(persistence persistence.IZonePersistence) *ZonePersistenceFixture {
	c := ZonePersistenceFixture{}

	c.ZONE1 = &data.ZoneV1{
		Id:        "1",
		MapId:     "map1",
		OrgId:     "org_1",
		Name:      "Zone 1",
		PositionX: 10,
		PositionY: 10,
		Width:     100,
		Height:    100,
		Type:      "type1",
		CreatedAt: time.Now(),
	}

	c.ZONE2 = &data.ZoneV1{
		Id:        "2",
		MapId:     "map1",
		OrgId:     "org_1",
		Name:      "Zone 2",
		PositionX: 20,
		PositionY: 20,
		Width:     200,
		Height:    200,
		Type:      "type2",
		CreatedAt: time.Now(),
	}

	c.ZONE3 = &data.ZoneV1{
		Id:        "3",
		MapId:     "map2",
		OrgId:     "org_2",
		Name:      "Zone 3",
		PositionX: 30,
		PositionY: 30,
		Width:     300,
		Height:    300,
		Type:      "type1",
		CreatedAt: time.Now(),
	}

	c.persistence = persistence
	return &c
}

func (c *ZonePersistenceFixture) testCreateZones(t *testing.T) {
	// Create the first zone
	z, err := c.persistence.Create(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, *c.ZONE1)
	assert.Nil(t, err)
	assert.NotNil(t, z)
	assert.Equal(t, c.ZONE1.Name, z.Name)
	assert.Equal(t, c.ZONE1.OrgId, z.OrgId)

	// Create the second zone
	z, err = c.persistence.Create(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, *c.ZONE2)
	assert.Nil(t, err)
	assert.NotNil(t, z)
	assert.Equal(t, c.ZONE2.Name, z.Name)
	assert.Equal(t, c.ZONE2.OrgId, z.OrgId)

	// Create the third zone
	z, err = c.persistence.Create(context.Background(), cdata.RequestContextV1{OrgId: "org_2"}, *c.ZONE3)
	assert.Nil(t, err)
	assert.NotNil(t, z)
	assert.Equal(t, c.ZONE3.Name, z.Name)
	assert.Equal(t, c.ZONE3.OrgId, z.OrgId)
}

func (c *ZonePersistenceFixture) TestCrudOperations(t *testing.T) {
	var zone1 data.ZoneV1

	// Create items
	c.testCreateZones(t)

	// Get all zones
	page, err := c.persistence.GetPageByFilter(context.Background(), cdata.RequestContextV1{},
		*cquery.NewEmptyFilterParams(), *cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 3)
	zone1 = page.Data[0]

	// Update the zone
	zone1.Name = "Updated Zone 1"
	z, err := c.persistence.Update(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, zone1)
	assert.Nil(t, err)
	assert.NotNil(t, z)
	assert.Equal(t, zone1.Id, z.Id)
	assert.Equal(t, "Updated Zone 1", z.Name)

	// Delete the zone
	z, err = c.persistence.DeleteById(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, zone1.Id)
	assert.Nil(t, err)
	assert.NotNil(t, z)
	assert.Equal(t, zone1.Id, z.Id)

	// Try to get deleted zone
	z, err = c.persistence.GetOneById(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, zone1.Id)
	assert.Nil(t, err)
	assert.Equal(t, data.ZoneV1{}, z)
}

func (c *ZonePersistenceFixture) TestGetWithFilters(t *testing.T) {
	// Create items
	c.testCreateZones(t)

	// Filter by id
	page, err := c.persistence.GetPageByFilter(context.Background(), cdata.RequestContextV1{},
		*cquery.NewFilterParamsFromTuples("id", "1"),
		*cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.Len(t, page.Data, 1)

	// Filter by org_id
	page, err = c.persistence.GetPageByFilter(context.Background(), cdata.RequestContextV1{OrgId: "org_1"},
		*cquery.NewFilterParamsFromTuples("org_id", "org_1"),
		*cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.Len(t, page.Data, 2)

	// Filter by map_id
	page, err = c.persistence.GetPageByFilter(context.Background(), cdata.RequestContextV1{},
		*cquery.NewFilterParamsFromTuples("map_id", "map1"),
		*cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.Len(t, page.Data, 2)

	// Filter by type
	page, err = c.persistence.GetPageByFilter(context.Background(), cdata.RequestContextV1{},
		*cquery.NewFilterParamsFromTuples("type", "type1"),
		*cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.Len(t, page.Data, 2)
}
