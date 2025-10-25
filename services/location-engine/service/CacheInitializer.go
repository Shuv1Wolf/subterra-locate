package service

import (
	"context"
	"encoding/json"
	"time"

	natsEvents "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/utils"
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

			c.beaconStateStore.Upsert(&utils.BeaconState{
				OrgID:      beacon.OrgId,
				MapID:      beacon.MapId,
				BeaconID:   beacon.Id,
				BeaconName: beacon.Label,
				X:          beacon.X,
				Y:          beacon.Y,
				Z:          beacon.Z,
				Info: map[string]string{
					"source": "system",
				},
				UpdatedAt: time.Now(),
			})
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

	if !b.Enabled {
		c.beaconStateStore.Upsert(&utils.BeaconState{
			OrgID:      b.OrgId,
			MapID:      b.MapId,
			BeaconID:   b.Id,
			BeaconName: b.Label,
			X:          0,
			Y:          0,
			Z:          0,
			Info: map[string]string{
				"source": "system",
			},
			UpdatedAt: time.Now(),
		})
	} else {
		c.beaconStateStore.Upsert(&utils.BeaconState{
			OrgID:      b.OrgId,
			MapID:      b.MapId,
			BeaconID:   b.Id,
			BeaconName: b.Label,
			X:          b.X,
			Y:          b.Y,
			Z:          b.Z,
			Info: map[string]string{
				"source": "system",
			},
			UpdatedAt: time.Now(),
		})
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

	if be, ok := c.beaconsMap[event.Id]; ok {
		c.beaconStateStore.Upsert(&utils.BeaconState{
			OrgID:      be.OrgId,
			MapID:      be.MapId,
			BeaconID:   be.Id,
			BeaconName: be.Label,
			X:          0,
			Y:          0,
			Z:          0,
			Info: map[string]string{
				"source": "system",
			},
			UpdatedAt: time.Now(),
		})
	}

	delete(c.beaconsMap, event.Id)
	return nil
}

// -------------------- Device cache --------------------
func (c *LocationEngineService) initDeviceCache() {
	filter := query.NewFilterParamsFromTuples("enabled", true)
	limit := int64(100)
	skip := int64(0)

	c.mu.Lock()
	defer c.mu.Unlock()

	for {
		page := query.NewPagingParams(skip, limit, false)
		res, err := c.deviceAdmin.GetDevices(context.Background(), *filter, *page)
		if err != nil {
			c.Logger.Error(context.Background(), err, "Failed to get devices from device admin service")
			return
		}

		if len(res.Data) == 0 {
			break
		}

		for _, device := range res.Data {
			c.deviceMap[device.Id] = &device
		}

		if int64(len(res.Data)) < limit {
			break
		}

		skip += limit
	}

	c.Logger.Info(context.Background(), "Devices stored in cache")
}

func (c *LocationEngineService) deviceChangedEvent(ctx context.Context, msg string) error {
	var event natsEvents.DeviceChangedEvent
	err := json.Unmarshal([]byte(msg), &event)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to deserialize message")
	}

	d, err := c.deviceAdmin.GetDeviceById(context.Background(), event.Id)
	if err != nil {
		c.Logger.Error(context.Background(), err, "Failed to get device from device admin service")
		return err
	}

	if !d.Enabled {
		c.deviceStateStore.Upsert(&utils.DeviceState{
			OrgID:      d.OrgId,
			MapID:      "",
			DeviceID:   d.Id,
			DeviceName: d.Name,
			X:          0,
			Y:          0,
			Z:          0,
			Info: map[string]string{
				"source": "ble",
			},
			UpdatedAt: time.Now(),
		})
	}

	c.deviceMap[event.Id] = d
	return nil
}

func (c *LocationEngineService) deviceDeletedEvent(ctx context.Context, msg string) error {
	var event natsEvents.DeviceChangedEvent
	err := json.Unmarshal([]byte(msg), &event)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to deserialize message")
	}

	if dev, ok := c.deviceMap[event.Id]; ok {
		c.deviceStateStore.Upsert(&utils.DeviceState{
			OrgID:      dev.OrgId,
			MapID:      "",
			DeviceID:   dev.Id,
			DeviceName: dev.Name,
			X:          0,
			Y:          0,
			Z:          0,
			Info: map[string]string{
				"source": "ble",
			},
			UpdatedAt: time.Now(),
		})
	}

	delete(c.deviceMap, event.Id)
	return nil
}
