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

type Map2dPersistenceFixture struct {
	MAP1        *data.Map2dV1
	MAP2        *data.Map2dV1
	MAP3        *data.Map2dV1
	persistence persistence.IMap2dPersistence
}

func NewMap2dPersistenceFixture(persistence persistence.IMap2dPersistence) *Map2dPersistenceFixture {
	c := Map2dPersistenceFixture{}

	c.MAP1 = &data.Map2dV1{
		Id:        "1",
		Name:      "TestMap1",
		SVG:       "<svg>map1</svg>",
		ScaleX:    1.0,
		ScaleY:    1.0,
		CreatedAt: time.Now(),
		OrgId:     "org_1",
		Width:     100,
		Height:    100,
		Level:     1,
	}

	c.MAP2 = &data.Map2dV1{
		Id:        "2",
		Name:      "TestMap2",
		SVG:       "<svg>map2</svg>",
		ScaleX:    2.0,
		ScaleY:    2.0,
		CreatedAt: time.Now(),
		OrgId:     "org_1",
		Width:     200,
		Height:    200,
		Level:     2,
	}

	c.MAP3 = &data.Map2dV1{
		Id:        "3",
		Name:      "TestMap3",
		SVG:       "<svg>map3</svg>",
		ScaleX:    3.0,
		ScaleY:    3.0,
		CreatedAt: time.Now(),
		OrgId:     "org_2",
		Width:     300,
		Height:    300,
		Level:     3,
	}

	c.persistence = persistence
	return &c
}

func (c *Map2dPersistenceFixture) testCreateMaps(t *testing.T) {
	// Create the first map
	m, err := c.persistence.Create(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, *c.MAP1)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, c.MAP1.Name, m.Name)
	assert.Equal(t, c.MAP1.OrgId, m.OrgId)

	// Create the second map
	m, err = c.persistence.Create(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, *c.MAP2)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, c.MAP2.Name, m.Name)
	assert.Equal(t, c.MAP2.OrgId, m.OrgId)

	// Create the third map
	m, err = c.persistence.Create(context.Background(), cdata.RequestContextV1{OrgId: "org_2"}, *c.MAP3)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, c.MAP3.Name, m.Name)
	assert.Equal(t, c.MAP3.OrgId, m.OrgId)
}

func (c *Map2dPersistenceFixture) TestCrudOperations(t *testing.T) {
	var map1 data.Map2dV1

	// Create items
	c.testCreateMaps(t)

	// Get all maps
	page, err := c.persistence.GetPageByFilter(context.Background(), cdata.RequestContextV1{},
		*cquery.NewEmptyFilterParams(), *cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 3)
	map1 = page.Data[0]

	// Update the map
	map1.Name = "Updated Map1"
	m, err := c.persistence.Update(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, map1)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, map1.Id, m.Id)
	assert.Equal(t, "Updated Map1", m.Name)

	// Delete the map
	m, err = c.persistence.DeleteById(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, map1.Id)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, map1.Id, m.Id)

	// Try to get deleted map
	m, err = c.persistence.GetOneById(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, map1.Id)
	assert.Nil(t, err)
	assert.Equal(t, data.Map2dV1{}, m)
}

func (c *Map2dPersistenceFixture) TestGetWithFilters(t *testing.T) {
	// Create items
	c.testCreateMaps(t)

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
}
