package utils

import (
	"sync"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/protos"
)

type changeZone struct {
	orgID string
	Ev    *protos.MonitorZoneStreamEventV1_ZoneEventV1
}

type zoneSubscriber struct {
	ch      chan changeZone
	closeCh chan struct{}
}

type orgBusZones struct {
	mu   sync.RWMutex
	subs map[*zoneSubscriber]struct{}
}

func newOrgBusZones() *orgBusZones {
	return &orgBusZones{subs: map[*zoneSubscriber]struct{}{}}
}

func (b *orgBusZones) subscribe() *zoneSubscriber {
	sub := &zoneSubscriber{
		ch:      make(chan changeZone, 64),
		closeCh: make(chan struct{}),
	}
	b.mu.Lock()
	b.subs[sub] = struct{}{}
	b.mu.Unlock()
	return sub
}

func (b *orgBusZones) unsubscribe(sub *zoneSubscriber) {
	b.mu.Lock()
	if _, ok := b.subs[sub]; ok {
		delete(b.subs, sub)
		close(sub.closeCh)
		close(sub.ch)
	}
	b.mu.Unlock()
}

func (b *orgBusZones) publish(event changeZone) {
	b.mu.RLock()
	for sub := range b.subs {
		s := sub
		go func() {
			select {
			case <-s.closeCh:
				return
			case s.ch <- event:
			}
		}()
	}
	b.mu.RUnlock()
}

type ZoneStateEntry struct {
	Zone *data1.ZoneV1
	Info map[string]string
}

type ZoneStateStore struct {
	mu   sync.RWMutex
	data map[string]map[string]*ZoneStateEntry
	bus  map[string]*orgBusZones
}

func NewZoneStateStore() *ZoneStateStore {
	return &ZoneStateStore{
		data: map[string]map[string]*ZoneStateEntry{},
		bus:  map[string]*orgBusZones{},
	}
}

func (s *ZoneStateStore) ensureBus(orgID string) *orgBusZones {
	if _, ok := s.bus[orgID]; !ok {
		s.bus[orgID] = newOrgBusZones()
	}
	return s.bus[orgID]
}

func (s *ZoneStateStore) Upsert(zone *data1.ZoneV1) {
	s.mu.Lock()

	if _, ok := s.data[zone.OrgId]; !ok {
		s.data[zone.OrgId] = map[string]*ZoneStateEntry{}
	}

	var info map[string]string
	if entry, ok := s.data[zone.OrgId][zone.Id]; ok {
		info = entry.Info
	} else {
		info = map[string]string{}
	}

	s.data[zone.OrgId][zone.Id] = &ZoneStateEntry{
		Zone: zone,
		Info: info,
	}

	b := s.ensureBus(zone.OrgId)
	s.mu.Unlock()

	b.publish(changeZone{
		orgID: zone.OrgId,
		Ev: &protos.MonitorZoneStreamEventV1_ZoneEventV1{
			ZoneId:    zone.Id,
			ZoneName:  zone.Name,
			MapId:     zone.MapId,
			OrgId:     zone.OrgId,
			PositionX: float32(zone.PositionX),
			PositionY: float32(zone.PositionY),
			Width:     float32(zone.Width),
			Height:    float32(zone.Height),
			Type:      string(zone.Type),
			Info:      info,
			Color:     zone.Color,
			Deleted:   false,
		},
	})
}

func (s *ZoneStateStore) UpdateState(orgID, zoneID string, info map[string]string) {
	s.mu.Lock()
	zones, ok := s.data[orgID]
	if !ok {
		s.mu.Unlock()
		return
	}
	entry, ok := zones[zoneID]
	if !ok {
		s.mu.Unlock()
		return
	}

	entry.Info = info
	zone := entry.Zone
	b := s.ensureBus(orgID)
	s.mu.Unlock()

	b.publish(changeZone{
		orgID: zone.OrgId,
		Ev: &protos.MonitorZoneStreamEventV1_ZoneEventV1{
			ZoneId:    zone.Id,
			ZoneName:  zone.Name,
			MapId:     zone.MapId,
			OrgId:     zone.OrgId,
			PositionX: float32(zone.PositionX),
			PositionY: float32(zone.PositionY),
			Width:     float32(zone.Width),
			Height:    float32(zone.Height),
			Type:      string(zone.Type),
			Info:      info,
			Color:     zone.Color,
			Deleted:   false,
		},
	})
}

func (s *ZoneStateStore) Delete(zone *data1.ZoneV1) {
	s.mu.RLock()
	b := s.bus[zone.OrgId]
	s.mu.RUnlock()

	if b != nil {
		b.publish(changeZone{
			orgID: zone.OrgId,
			Ev: &protos.MonitorZoneStreamEventV1_ZoneEventV1{
				ZoneId:    zone.Id,
				ZoneName:  zone.Name,
				MapId:     zone.MapId,
				OrgId:     zone.OrgId,
				PositionX: float32(zone.PositionX),
				PositionY: float32(zone.PositionY),
				Width:     float32(zone.Width),
				Height:    float32(zone.Height),
				Type:      string(zone.Type),
				Info:      map[string]string{},
				Color:     zone.Color,
				Deleted:   true,
			},
		})
	}

	s.mu.Lock()
	if zones, ok := s.data[zone.OrgId]; ok {
		delete(zones, zone.Id)
	}
	s.mu.Unlock()
}

func (s *ZoneStateStore) Snapshot(orgID string) []*protos.MonitorZoneStreamEventV1_ZoneEventV1 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	zs := s.data[orgID]
	if zs == nil {
		return nil
	}

	res := make([]*protos.MonitorZoneStreamEventV1_ZoneEventV1, 0, len(zs))
	for _, z := range zs {
		res = append(res, &protos.MonitorZoneStreamEventV1_ZoneEventV1{
			ZoneId:    z.Zone.Id,
			ZoneName:  z.Zone.Name,
			MapId:     z.Zone.MapId,
			OrgId:     z.Zone.OrgId,
			PositionX: float32(z.Zone.PositionX),
			PositionY: float32(z.Zone.PositionY),
			Width:     float32(z.Zone.Width),
			Height:    float32(z.Zone.Height),
			Type:      string(z.Zone.Type),
			Color:     z.Zone.Color,
			Info:      z.Info,
		})
	}
	return res
}

func (s *ZoneStateStore) GetZones(orgID string) []*data1.ZoneV1 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	zones := s.data[orgID]
	res := make([]*data1.ZoneV1, 0, len(zones))
	for _, z := range zones {
		res = append(res, z.Zone)
	}
	return res
}

type ZoneSubscription struct {
	bus *orgBusZones
	sub *zoneSubscriber
}

func (s *ZoneStateStore) Subscribe(orgID string) *ZoneSubscription {
	s.mu.Lock()
	b := s.ensureBus(orgID)
	sub := b.subscribe()
	s.mu.Unlock()
	return &ZoneSubscription{bus: b, sub: sub}
}

func (s *ZoneSubscription) C() <-chan changeZone { return s.sub.ch }

func (s *ZoneSubscription) Close() { s.bus.unsubscribe(s.sub) }
