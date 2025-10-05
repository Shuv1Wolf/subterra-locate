package utils

import (
	"sync"
	"time"

	protos "github.com/Shuv1Wolf/subterra-locate/services/location-engine/protos"
)

type BeaconState struct {
	OrgID      string
	MapID      string
	BeaconID   string
	BeaconName string
	X, Y, Z    float32
	Info       map[string]string
	UpdatedAt  time.Time
}

type changeBeacon struct {
	orgID string
	Ev    *protos.MonitorBeaconLocationStreamEventV1_LocationEventV1
}

type orgBusBeacon struct {
	mu   sync.RWMutex
	subs map[chan changeBeacon]struct{}
}

func newOrgBusBeacon() *orgBusBeacon { return &orgBusBeacon{subs: map[chan changeBeacon]struct{}{}} }

func (b *orgBusBeacon) subscribe() chan changeBeacon {
	ch := make(chan changeBeacon, 256)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *orgBusBeacon) unsubscribe(ch chan changeBeacon) {
	b.mu.Lock()
	delete(b.subs, ch)
	b.mu.Unlock()
	close(ch)
}

func (b *orgBusBeacon) publish(c changeBeacon) {
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

type BeaconStateStore struct {
	mu sync.RWMutex
	// orgID -> beaconID -> state
	byOrg map[string]map[string]*BeaconState
	// orgID -> bus
	bus map[string]*orgBusBeacon
}

func NewBeaconStateStore() *BeaconStateStore {
	return &BeaconStateStore{
		byOrg: map[string]map[string]*BeaconState{},
		bus:   map[string]*orgBusBeacon{},
	}
}

func (s *BeaconStateStore) Upsert(ev *BeaconState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byOrg[ev.OrgID]; !ok {
		s.byOrg[ev.OrgID] = map[string]*BeaconState{}
	}
	s.byOrg[ev.OrgID][ev.BeaconID] = ev

	if _, ok := s.bus[ev.OrgID]; !ok {
		s.bus[ev.OrgID] = newOrgBusBeacon()
	}
	s.bus[ev.OrgID].publish(changeBeacon{
		orgID: ev.OrgID,
		Ev: &protos.MonitorBeaconLocationStreamEventV1_LocationEventV1{
			BeaconId:   ev.BeaconID,
			BeaconName: ev.BeaconName,
			MapId:      ev.MapID,
			X:          ev.X,
			Y:          ev.Y,
			Z:          ev.Z,
			Info:       ev.Info,
		},
	})
}

func (s *BeaconStateStore) Snapshot(orgID string) []*protos.MonitorBeaconLocationStreamEventV1_LocationEventV1 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*protos.MonitorBeaconLocationStreamEventV1_LocationEventV1, 0)
	orgMap := s.byOrg[orgID]
	if orgMap == nil {
		return res
	}
	for _, st := range orgMap {
		res = append(res, &protos.MonitorBeaconLocationStreamEventV1_LocationEventV1{
			BeaconId:   st.BeaconID,
			BeaconName: st.BeaconName,
			MapId:      st.MapID,
			X:          st.X, Y: st.Y, Z: st.Z,
			Info: st.Info,
		})
	}
	return res
}

func (s *BeaconStateStore) Subscribe(orgID string) *orgBusSubscriptionBeacon {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.bus[orgID]; !ok {
		s.bus[orgID] = newOrgBusBeacon()
	}
	ch := s.bus[orgID].subscribe()
	return &orgBusSubscriptionBeacon{bus: s.bus[orgID], ch: ch}
}

type orgBusSubscriptionBeacon struct {
	bus *orgBusBeacon
	ch  chan changeBeacon
}

func (sub *orgBusSubscriptionBeacon) C() <-chan changeBeacon { return sub.ch }

func (sub *orgBusSubscriptionBeacon) Close() { sub.bus.unsubscribe(sub.ch) }
