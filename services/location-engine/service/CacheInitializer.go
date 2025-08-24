package service

import (
	"context"
	"encoding/json"

	natsEvents "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

// -------------------- Beacons cache --------------------
func (c *LocationEngineService) initBeaconsCache() {
	filter := query.NewFilterParamsFromTuples("enabled", true)
	limit := int64(100)
	skip := int64(0)

	c.mu.Lock()
	defer c.mu.Unlock()

	for {
		page := query.NewPagingParams(skip, limit, false)
		res, err := c.beaconAdmin.GetBeacons(context.Background(), *filter, *page)
		if err != nil {
			c.Logger.Error(context.Background(), err, "Failed to get beacons from beacon admin service")
			return
		}

		if len(res.Data) == 0 {
			break
		}

		for _, beacon := range res.Data {
			c.beaconsMap[beacon.Id] = &beacon
		}

		if int64(len(res.Data)) < limit {
			break
		}

		skip += limit
	}

	c.Logger.Info(context.Background(), "Beacons stored in cache")
}

func (c *LocationEngineService) beaconChangedEvent(ctx context.Context, msg string) error {
	var event natsEvents.BeaconChangedEvent
	err := json.Unmarshal([]byte(msg), &event)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to deserialize message")
	}

	b, err := c.beaconAdmin.GetBeaconById(context.Background(), event.Id)
	if err != nil {
		c.Logger.Error(context.Background(), err, "Failed to get beacon from beacon admin service")
		return err
	}

	c.beaconsMap[event.Id] = b
	return nil
}

func (c *LocationEngineService) beaconDeletedEvent(ctx context.Context, msg string) error {
	var event natsEvents.BeaconChangedEvent
	err := json.Unmarshal([]byte(msg), &event)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to deserialize message")
	}

	delete(c.beaconsMap, event.Id)
	return nil
}

// -------------------- Device cache --------------------
