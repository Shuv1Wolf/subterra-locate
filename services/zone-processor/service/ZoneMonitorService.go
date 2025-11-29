package service

import (
	"time"

	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/protos"
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/utils"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

const heartbeatInterval = 2 * time.Second

type ZoneMonitorService struct {
	protos.UnimplementedZoneMonitorServer
	stateStore *utils.ZoneStateStore
	logger     *clog.CompositeLogger
}

func NewZoneMonitorService(stateStore *utils.ZoneStateStore, logger *clog.CompositeLogger) *ZoneMonitorService {
	return &ZoneMonitorService{
		stateStore: stateStore,
		logger:     logger,
	}
}

func (s *ZoneMonitorService) MonitorZoneV1(
	in *protos.MonitorZoneRequestV1,
	stream protos.ZoneMonitor_MonitorZoneV1Server,
) error {
	ctx := stream.Context()
	orgID := in.GetOrgId()
	mapID := in.GetMapId()

	initial := s.stateStore.Snapshot(orgID)
	resp := make([]*protos.MonitorZoneStreamEventV1_ZoneEventV1, 0, len(initial))
	for _, ev := range initial {
		if mapID != "" && ev.GetMapId() != mapID {
			continue
		}

		resp = append(resp, ev)
	}
	if err := stream.Send(&protos.MonitorZoneStreamEventV1{Event: resp}); err != nil {
		return err
	}

	sub := s.stateStore.Subscribe(orgID)
	defer sub.Close()

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	pending := map[string]*protos.MonitorZoneStreamEventV1_ZoneEventV1{}

	flush := func() error {
		if len(pending) == 0 {
			return stream.Send(&protos.MonitorZoneStreamEventV1{
				Event: nil,
			})
		}
		batch := make([]*protos.MonitorZoneStreamEventV1_ZoneEventV1, 0, len(pending))
		for _, ev := range pending {
			batch = append(batch, ev)
		}
		pending = map[string]*protos.MonitorZoneStreamEventV1_ZoneEventV1{}
		return stream.Send(&protos.MonitorZoneStreamEventV1{Event: batch})
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case c := <-sub.C():
			if mapID != "" && c.Ev.GetMapId() != mapID {
				continue
			}

			pending[c.Ev.GetZoneId()] = c.Ev

		case <-ticker.C:
			if err := flush(); err != nil {
				return err
			}
		}
	}
}
