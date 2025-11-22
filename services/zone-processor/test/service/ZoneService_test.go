package test_service

import (
	"context"
	"testing"
	"time"

	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	data "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/persistence"
	logic "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/service"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

type ZoneServiceTest struct {
	ZONE1    *data.ZoneV1
	ZONE2    *data.ZoneV1
	service *logic.ZoneService
}

func newZoneServiceTest() *ZoneServiceTest {
	zonePersistence := persistence.NewZoneMemoryPersistence()
	zonePersistence.Configure(context.Background(), cconf.NewEmptyConfigParams())

	srv := logic.NewZoneService()
	srv.Configure(context.Background(), cconf.NewEmptyConfigParams())

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("zone-processor", "persistence", "memory", "default", "1.0"), zonePersistence,
	)
	srv.SetReferences(context.Background(), references)

	return &ZoneServiceTest{
		ZONE1: &data.ZoneV1{
			Id:        "1",
			MapId:     "map_1",
			OrgId:     "org_1",
			Name:      "TestZone1",
			PositionX: 10,
			PositionY: 10,
			Width:     100,
			Height:    100,
			Type:      "polygon",
			MaxDevice: 10,
			CreatedAt: time.Now(),
		},
		ZONE2: &data.ZoneV1{
			Id:        "2",
			MapId:     "map_1",
			OrgId:     "org_1",
			Name:      "TestZone2",
			PositionX: 20,
			PositionY: 20,
			Width:     200,
			Height:    200,
			Type:      "circle",
			MaxDevice: 20,
			CreatedAt: time.Now(),
		},
		service: srv,
	}
}

func (c *ZoneServiceTest) TestZoneCrudOperations(t *testing.T) {
	// Create one zone
	z, err := c.service.CreateZone(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, *c.ZONE1)
	assert.Nil(t, err)
	assert.NotNil(t, z)
	assert.Equal(t, c.ZONE1.Id, z.Id)
	assert.Equal(t, c.ZONE1.Name, z.Name)

	// Create another zone
	z, err = c.service.CreateZone(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, *c.ZONE2)
	assert.Nil(t, err)
	assert.NotNil(t, z)
	assert.Equal(t, c.ZONE2.Id, z.Id)
	assert.Equal(t, c.ZONE2.Name, z.Name)

	// Get all zones
	page, err := c.service.GetZones(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, *cquery.NewEmptyFilterParams(), *cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 2)

	// Update the zone
	z.Name = "Updated Zone2"
	updatedZone, err := c.service.UpdateZone(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, z)
	assert.Nil(t, err)
	assert.NotNil(t, updatedZone)
	assert.Equal(t, z.Id, updatedZone.Id)
	assert.Equal(t, "Updated Zone2", updatedZone.Name)

	// Delete zone
	deletedZone, err := c.service.DeleteZoneById(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, z.Id)
	assert.Nil(t, err)
	assert.NotNil(t, deletedZone)
	assert.Equal(t, z.Id, deletedZone.Id)

	// Try to get deleted zone
	getZone, err := c.service.GetZoneById(context.Background(), cdata.RequestContextV1{OrgId: "org_1"}, z.Id)
	assert.Nil(t, err)
	assert.Equal(t, data.ZoneV1{}, getZone)
}

func TestZoneService(t *testing.T) {
	c := newZoneServiceTest()
	t.Run("Zone CRUD Operations", c.TestZoneCrudOperations)
}
