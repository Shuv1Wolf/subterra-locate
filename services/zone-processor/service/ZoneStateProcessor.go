package service

import (
	"strconv"
	"strings"
	"sync"

	"github.com/Shuv1Wolf/subterra-locate/services/common/nats/events"
	data1 "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/utils"
)

type ZoneStateProcessor struct {
	store *utils.ZoneStateStore
	// deviceZones maps orgID -> deviceID -> zoneID
	deviceZones map[string]map[string]string
	mu          sync.RWMutex
}

func NewZoneStateProcessor(store *utils.ZoneStateStore) *ZoneStateProcessor {
	return &ZoneStateProcessor{
		store:       store,
		deviceZones: make(map[string]map[string]string),
	}
}

func (p *ZoneStateProcessor) ProcessPosition(event events.DevicePositioningEventV1) {
	p.mu.Lock()
	defer p.mu.Unlock()

	orgID := event.OrgId
	deviceID := event.DeviceId
	mapID := event.MapId

	// Ensure org map exists
	if _, ok := p.deviceZones[orgID]; !ok {
		p.deviceZones[orgID] = make(map[string]string)
	}

	// Get current zone for device
	currentZoneID, hasCurrentZone := p.deviceZones[orgID][deviceID]

	// Find new zone
	var newZoneID string
	zones := p.store.GetZones(orgID) // This acquires store lock

	for _, zone := range zones {
		// Check map ID match
		if zone.MapId != mapID {
			continue
		}

		// Check if point in rectangle
		if event.X >= zone.PositionX && event.X <= zone.PositionX+zone.Width &&
			event.Y >= zone.PositionY && event.Y <= zone.PositionY+zone.Height {
			newZoneID = zone.Id
			break // Assume device can be in only one zone (or take first one)
		}
	}

	if newZoneID == currentZoneID {
		return // No change
	}

	// Process Exit
	if hasCurrentZone && currentZoneID != "" {
		p.updateZoneState(orgID, currentZoneID, deviceID, false)
	}

	// Process Enter
	if newZoneID != "" {
		p.updateZoneState(orgID, newZoneID, deviceID, true)
		p.deviceZones[orgID][deviceID] = newZoneID
	} else {
		// Device moved out of any zone
		delete(p.deviceZones[orgID], deviceID)
	}
}

func (p *ZoneStateProcessor) updateZoneState(orgID, zoneID, deviceID string, entered bool) {
	// Need to calculate new state for the zone.
	// We can iterate p.deviceZones[orgID] to find all devices in zoneID.

	devices := []string{}
	for dID, zID := range p.deviceZones[orgID] {
		if zID == zoneID && dID != deviceID { // Don't include current device yet if we are updating
			devices = append(devices, dID)
		}
	}

	if entered {
		devices = append(devices, deviceID)
	}

	count := len(devices)

	info := map[string]string{
		"count":   strconv.Itoa(count),
		"devices": strings.Join(devices, ","),
	}

	if entered {
		info["last_entered"] = deviceID
		info["last_exited"] = ""
	} else {
		info["last_entered"] = ""
		info["last_exited"] = deviceID
	}

	p.store.UpdateState(orgID, zoneID, info)
}

func (p *ZoneStateProcessor) HandleZoneUpdate(zone *data1.ZoneV1) {
	p.mu.Lock()
	defer p.mu.Unlock()

	orgID := zone.OrgId
	zoneID := zone.Id

	if _, ok := p.deviceZones[orgID]; !ok {
		return
	}

	// Remove all devices from this zone in our cache
	// Note: Iterating map while deleting is safe in Go
	for dID, zID := range p.deviceZones[orgID] {
		if zID == zoneID {
			delete(p.deviceZones[orgID], dID)
		}
	}

	// Update store to empty state
	p.store.UpdateState(orgID, zoneID, map[string]string{
		"count":        "0",
		"devices":      "",
		"last_entered": "",
		"last_exited":  "",
	})
}

func (p *ZoneStateProcessor) HandleZoneDelete(zone *data1.ZoneV1) {
	p.mu.Lock()
	defer p.mu.Unlock()

	orgID := zone.OrgId
	zoneID := zone.Id

	if _, ok := p.deviceZones[orgID]; !ok {
		return
	}

	for dID, zID := range p.deviceZones[orgID] {
		if zID == zoneID {
			delete(p.deviceZones[orgID], dID)
		}
	}
}

func (p *ZoneStateProcessor) HandleZoneAdd(zone *data1.ZoneV1) {
	p.HandleZoneUpdate(zone)
}
