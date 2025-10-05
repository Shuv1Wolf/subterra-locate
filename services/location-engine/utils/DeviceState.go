package utils

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

type changeDevice struct {
	orgID string
	Ev    *protos.MonitorDeviceLocationStreamEventV1_LocationEventV1
}

type orgBusDevice struct {
	mu   sync.RWMutex
	subs map[chan changeDevice]struct{}
}

func newOrgBusDevice() *orgBusDevice { return &orgBusDevice{subs: map[chan changeDevice]struct{}{}} }

func (b *orgBusDevice) subscribe() chan changeDevice {
	ch := make(chan changeDevice, 256)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *orgBusDevice) unsubscribe(ch chan changeDevice) {
	b.mu.Lock()
	delete(b.subs, ch)
	b.mu.Unlock()
	close(ch)
}

func (b *orgBusDevice) publish(c changeDevice) {
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

type DeviceStateStore struct {
	mu sync.RWMutex
	// orgID -> deviceID -> state
	byOrg map[string]map[string]*DeviceState
	// orgID -> bus
	bus map[string]*orgBusDevice
}

func NewDeviceStateStore() *DeviceStateStore {
	return &DeviceStateStore{
		byOrg: map[string]map[string]*DeviceState{},
		bus:   map[string]*orgBusDevice{},
	}
}

func (s *DeviceStateStore) Upsert(ev *DeviceState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byOrg[ev.OrgID]; !ok {
		s.byOrg[ev.OrgID] = map[string]*DeviceState{}
	}
	s.byOrg[ev.OrgID][ev.DeviceID] = ev

	if _, ok := s.bus[ev.OrgID]; !ok {
		s.bus[ev.OrgID] = newOrgBusDevice()
	}
	s.bus[ev.OrgID].publish(changeDevice{
		orgID: ev.OrgID,
		Ev: &protos.MonitorDeviceLocationStreamEventV1_LocationEventV1{
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

func (s *DeviceStateStore) Snapshot(orgID string) []*protos.MonitorDeviceLocationStreamEventV1_LocationEventV1 {
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

func (s *DeviceStateStore) Subscribe(orgID string) *orgBusSubscriptionDevice {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.bus[orgID]; !ok {
		s.bus[orgID] = newOrgBusDevice()
	}
	ch := s.bus[orgID].subscribe()
	return &orgBusSubscriptionDevice{bus: s.bus[orgID], ch: ch}
}

type orgBusSubscriptionDevice struct {
	bus *orgBusDevice
	ch  chan changeDevice
}

func (sub *orgBusSubscriptionDevice) C() <-chan changeDevice { return sub.ch }

func (sub *orgBusSubscriptionDevice) Close() { sub.bus.unsubscribe(sub.ch) }
