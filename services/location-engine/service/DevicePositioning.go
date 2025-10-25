package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	natsEvents "github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/utils"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
)

// EstimateXYZ estimates the device position (x,y,z) using multilateration.
// n - path-loss exponent (≈2.0 office, 2.2–3.5 concrete/warehouse).
// defaultTxPower - fallback txPower (often -59 dBm) if beacon Txpower=0.
func (c *LocationEngineService) EstimateXYZ(
	ctx context.Context,
	event natsEvents.DeviceDetectedBLERawEventV1,
	n float64,
	defaultTxPower int,
) (x, y, z float64, mapId string, err error) {

	mapIdCounts := make(map[string]int)
	for _, e := range event.Beacons {
		c.mu.Lock()
		b, ok := c.beaconsMap[e.Id]
		c.mu.Unlock()
		if ok && b.Enabled {
			mapIdCounts[b.MapId]++
		}
	}

	var mostFrequentMapId string
	maxCount := 0
	for mapId, count := range mapIdCounts {
		if count > maxCount {
			maxCount = count
			mostFrequentMapId = mapId
		}
	}

	var data []utils.Obs
	for _, e := range event.Beacons {

		c.mu.Lock()
		b, ok := c.beaconsMap[e.Id]
		c.mu.Unlock()

		if !ok || !b.Enabled {
			continue
		}

		txp := e.Txpower
		if txp == 0 {
			txp = defaultTxPower
		}
		d := utils.RssiToDistance(float64(e.Rssi), float64(txp), n)
		if !utils.IsFinite(d) || d <= 0 {
			continue
		}
		data = append(data, utils.Obs{
			X: float64(b.X), Y: float64(b.Y), Z: float64(b.Z),
			D: d,
			W: 1.0 / (d * d), // weight ~ inverse square of distance
		})
	}

	if len(data) < 3 {
		return 0, 0, 0, "", fmt.Errorf("need >=3 beacons, got %d", len(data))
	}

	utils.NormalizeWeights(data)

	// If all beacons are on approximately the same Z plane, solve in 2D.
	if planeZ, ok := utils.CommonZ(data, 0.25); ok {
		xx, yy, _, e := utils.GaussNewton(data, true, planeZ)
		return xx, yy, planeZ, mostFrequentMapId, e
	}

	xx, yy, zz, e := utils.GaussNewton(data, false, 0)
	return xx, yy, zz, mostFrequentMapId, e
}

func (c *LocationEngineService) bleEventHandler(ctx context.Context, envelope *cqueues.MessageEnvelope) error {
	var event natsEvents.DeviceDetectedBLERawEventV1
	err := json.Unmarshal([]byte(envelope.GetMessageAsString()), &event)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to deserialize message")
	}

	c.mu.Lock()
	d, ok := c.deviceMap[event.DeviceId]
	c.mu.Unlock()
	if !ok || !d.Enabled {
		c.Logger.Warn(ctx, "Device not found: "+event.DeviceId)
		return nil
	}

	x, y, z, mapId, err := c.EstimateXYZ(ctx, event, 3, -59)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to estimate XYZ")
	}

	pos := &natsEvents.DevicePositioningEventV1{
		DeviceId: event.DeviceId, X: x, Y: y, Z: z, Time: time.Now(),
	}

	c.deviceStateStore.Upsert(&utils.DeviceState{
		OrgID:      d.OrgId,
		MapID:      mapId,
		DeviceID:   pos.DeviceId,
		DeviceName: d.Name,
		X:          float32(pos.X),
		Y:          float32(pos.Y),
		Z:          float32(pos.Z),
		Info: map[string]string{
			"source": "ble",
			"time":   pos.Time.Format(time.RFC3339),
		},
		UpdatedAt: time.Now(),
		Online:     true,
	})

	// Only send position if the device is online
	if d, ok := c.deviceMap[event.DeviceId]; ok && d.Enabled {
		err = c.devicePositionPublisher.SendDevicePosition(ctx, pos)
		if err != nil {
			c.Logger.Error(ctx, err, "Failed to send message")
		}
	}
	c.Logger.Debug(ctx, "Received message: "+envelope.GetMessageAsString())
	return nil
}
