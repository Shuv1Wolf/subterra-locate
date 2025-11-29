package utils

import (
	"sync"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/data/version1"
	protos "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/protos"
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

type ZoneStateStore struct {
	mu   sync.RWMutex
	data map[string]map[string]*data1.ZoneV1
	bus  map[string]*orgBusZones
}

func NewZoneStateStore() *ZoneStateStore {
	return &ZoneStateStore{
		data: map[string]map[string]*data1.ZoneV1{},
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
		s.data[zone.OrgId] = map[string]*data1.ZoneV1{}
	}
	s.data[zone.OrgId][zone.Id] = zone

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
			Info:      map[string]string{},
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
			ZoneId:    z.Id,
			ZoneName:  z.Name,
			MapId:     z.MapId,
			OrgId:     z.OrgId,
			PositionX: float32(z.PositionX),
			PositionY: float32(z.PositionY),
			Width:     float32(z.Width),
			Height:    float32(z.Height),
			Type:      string(z.Type),
			Info:      map[string]string{},
		})
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
