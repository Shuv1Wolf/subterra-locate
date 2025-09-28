package service

import (
	"sync"
	"time"

	protos "github.com/Shuv1Wolf/subterra-locate/services/location-engine/protos"
)

type DeviceState struct {
	OrgID      string
	MapID      string
	DeviceID   string
	DeviceName string
	X, Y, Z    float32
	Info       map[string]string
	UpdatedAt  time.Time
}

type change struct {
	orgID string
	ev    *protos.MonitorDeviceLocationStreamEventV1_LocationEventV1
}

type orgBus struct {
	mu   sync.RWMutex
	subs map[chan change]struct{}
}

func newOrgBus() *orgBus { return &orgBus{subs: map[chan change]struct{}{}} }

func (b *orgBus) subscribe() chan change {
	ch := make(chan change, 256)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *orgBus) unsubscribe(ch chan change) {
	b.mu.Lock()
	delete(b.subs, ch)
	b.mu.Unlock()
	close(ch)
}

func (b *orgBus) publish(c change) {
	b.mu.RLock()
	for ch := range b.subs {
		select {
		case ch <- c:
		default:
			// drop
		}
	}
	b.mu.RUnlock()
}

type StateStore struct {
	mu sync.RWMutex
	// orgID -> deviceID -> state
	byOrg map[string]map[string]*DeviceState
	// orgID -> bus
	bus map[string]*orgBus
}

func NewStateStore() *StateStore {
	return &StateStore{
		byOrg: map[string]map[string]*DeviceState{},
		bus:   map[string]*orgBus{},
	}
}

func (s *StateStore) upsert(ev *DeviceState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byOrg[ev.OrgID]; !ok {
		s.byOrg[ev.OrgID] = map[string]*DeviceState{}
	}
	s.byOrg[ev.OrgID][ev.DeviceID] = ev

	if _, ok := s.bus[ev.OrgID]; !ok {
		s.bus[ev.OrgID] = newOrgBus()
	}
	s.bus[ev.OrgID].publish(change{
		orgID: ev.OrgID,
		ev: &protos.MonitorDeviceLocationStreamEventV1_LocationEventV1{
			DeviceId:   ev.DeviceID,
			DeviceName: ev.DeviceName,
			MapId:      ev.MapID,
			X:          ev.X,
			Y:          ev.Y,
			Z:          ev.Z,
			Info:       ev.Info,
		},
	})
}

func (s *StateStore) snapshot(orgID string) []*protos.MonitorDeviceLocationStreamEventV1_LocationEventV1 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*protos.MonitorDeviceLocationStreamEventV1_LocationEventV1, 0)
	orgMap := s.byOrg[orgID]
	if orgMap == nil {
		return res
	}
	for _, st := range orgMap {
		res = append(res, &protos.MonitorDeviceLocationStreamEventV1_LocationEventV1{
			DeviceId:   st.DeviceID,
			DeviceName: st.DeviceName,
			MapId:      st.MapID,
			X:          st.X, Y: st.Y, Z: st.Z,
			Info: st.Info,
		})
	}
	return res
}

func (s *StateStore) subscribe(orgID string) *orgBusSubscription {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.bus[orgID]; !ok {
		s.bus[orgID] = newOrgBus()
	}
	ch := s.bus[orgID].subscribe()
	return &orgBusSubscription{bus: s.bus[orgID], ch: ch}
}

type orgBusSubscription struct {
	bus *orgBus
	ch  chan change
}

func (sub *orgBusSubscription) C() <-chan change { return sub.ch }

func (sub *orgBusSubscription) Close() { sub.bus.unsubscribe(sub.ch) }
