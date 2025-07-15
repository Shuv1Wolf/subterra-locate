package logic

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/persistence"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type BeaconsService struct {
	persistence persistence.IBeaconsPersistence
	commandSet  *BeaconsCommandSet
}

func NewBeaconsService() *BeaconsService {
	c := &BeaconsService{}
	return c
}

func (c *BeaconsService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	// Read configuration parameters here...
}

func (c *BeaconsService) GetCommandSet() *ccmd.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewBeaconsCommandSet(c)
	}
	return &c.commandSet.CommandSet
}

func (c *BeaconsService) SetReferences(ctx context.Context, references cref.IReferences) {
	locator := cref.NewDescriptor("beacon-admin", "persistence", "*", "*", "1.0")
	p, err := references.GetOneRequired(locator)
	if p != nil && err == nil {
		if _pers, ok := p.(persistence.IBeaconsPersistence); ok {
			c.persistence = _pers
			return
		}
	}
	panic(cref.NewReferenceError(ctx, locator))
}

func (c *BeaconsService) GetBeacons(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.BeaconV1], error) {
	return c.persistence.GetPageByFilter(ctx, filter, paging)
}

func (c *BeaconsService) GetBeaconById(ctx context.Context,
	beaconId string) (data.BeaconV1, error) {

	return c.persistence.GetOneById(ctx, beaconId)
}

func (c *BeaconsService) GetBeaconByUdi(ctx context.Context,
	beaconId string) (data.BeaconV1, error) {

	return c.persistence.GetOneByUdi(ctx, beaconId)
}

func (c *BeaconsService) CreateBeacon(ctx context.Context,
	beacon data.BeaconV1) (data.BeaconV1, error) {

	if beacon.Id == "" {
		beacon.Id = cdata.IdGenerator.NextLong()
	}

	if beacon.Type == "" {
		beacon.Type = data.Unknown
	}

	return c.persistence.Create(ctx, beacon)
}

func (c *BeaconsService) UpdateBeacon(ctx context.Context,
	beacon data.BeaconV1) (data.BeaconV1, error) {

	if beacon.Type == "" {
		beacon.Type = data.Unknown
	}

	return c.persistence.Update(ctx, beacon)
}

func (c *BeaconsService) DeleteBeaconById(ctx context.Context,
	beaconId string) (data.BeaconV1, error) {

	return c.persistence.DeleteById(ctx, beaconId)
}
