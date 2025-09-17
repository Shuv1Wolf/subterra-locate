package logic

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/persistence"
	"github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/publisher"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type MapService struct {
	persistence persistence.IMap2dPersistence
	commandSet  *MapCommandSet
	map2dEvents publisher.IPublisher
}

func NewMapService() *MapService {
	c := &MapService{}
	return c
}

func (c *MapService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	// Read configuration parameters here...
}

func (c *MapService) GetCommandSet() *ccmd.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewMapCommandSet(c)
	}
	return &c.commandSet.CommandSet
}

func (c *MapService) SetReferences(ctx context.Context, references cref.IReferences) {
	res, err := references.GetOneRequired(
		cref.NewDescriptor("geo-renderer", "persistence", "*", "map-2d", "1.0"),
	)
	if err != nil {
		panic(err)
	}
	c.persistence = res.(persistence.IMap2dPersistence)

	res = references.GetOneOptional(
		cref.NewDescriptor("geo-renderer", "publisher", "nats", "map-2d-events", "1.0"),
	)
	c.map2dEvents = res.(publisher.IPublisher)
}

func (c *MapService) GetMaps(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.Map2dV1], error) {
	return c.persistence.GetPageByFilter(ctx, filter, paging)
}

func (c *MapService) GetMapById(ctx context.Context,
	id string) (data.Map2dV1, error) {

	return c.persistence.GetOneById(ctx, id)
}

func (c *MapService) CreateMap(ctx context.Context,
	map2d data.Map2dV1) (data.Map2dV1, error) {
	b, err := c.persistence.Create(ctx, map2d)
	if err != nil {
		return b, err
	}

	if c.map2dEvents != nil {
		err = c.map2dEvents.SendMap2dCreatedEvent(ctx, b.Id)
		if err != nil {
			return b, err
		}
	}

	return b, nil
}

func (c *MapService) UpdateMap(ctx context.Context,
	map2d data.Map2dV1) (data.Map2dV1, error) {

	b, err := c.persistence.Update(ctx, map2d)
	if err != nil {
		return b, err
	}

	if c.map2dEvents != nil {
		err = c.map2dEvents.SendMap2dChangedEvent(ctx, b.Id)
		if err != nil {
			return b, err
		}
	}

	return b, err
}

func (c *MapService) DeleteMapById(ctx context.Context,
	id string) (data.Map2dV1, error) {

	b, err := c.persistence.DeleteById(ctx, id)
	if err != nil {
		return b, err
	}

	if c.map2dEvents != nil {
		err = c.map2dEvents.SendMap2dDeletedEvent(ctx, b.Id)
		if err != nil {
			return b, err
		}
	}

	return b, err
}
