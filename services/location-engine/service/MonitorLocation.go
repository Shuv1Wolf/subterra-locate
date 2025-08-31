package service

import (
	"time"

	protos "github.com/Shuv1Wolf/subterra-locate/services/location-engine/protos"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	grpc "google.golang.org/grpc"
)

const heartbeatInterval = 5 * time.Second

type MonitorLocation struct {
	protos.UnimplementedLocationMonitorServer
	state  *StateStore
	logger *clog.CompositeLogger
}

func NewMonitorLocation(state *StateStore, logger *clog.CompositeLogger) *MonitorLocation {
	return &MonitorLocation{state: state, logger: logger}
}

func (s *MonitorLocation) MonitorDeviceLocationV1(
	in *protos.MonitorDeviceLocationRequestV1,
	stream grpc.ServerStreamingServer[protos.MonitorDeviceLocationStreamEventV1],
) error {
	ctx := stream.Context()
	orgID := in.GetOrgId()

	filter := make(map[string]struct{}, len(in.GetDeviceId()))
	for _, id := range in.GetDeviceId() {
		filter[id] = struct{}{}
	}

	initial := s.state.snapshot(orgID, filter)
	if err := stream.Send(&protos.MonitorDeviceLocationStreamEventV1{Event: initial}); err != nil {
		return err
	}

	sub := s.state.subscribe(orgID)
	defer sub.Close()

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	pending := map[string]*protos.MonitorDeviceLocationStreamEventV1_LocationEventV1{}

	flush := func() error {
		if len(pending) == 0 {
			return stream.Send(&protos.MonitorDeviceLocationStreamEventV1{
				Event: nil,
			})
		}
		batch := make([]*protos.MonitorDeviceLocationStreamEventV1_LocationEventV1, 0, len(pending))
		for _, ev := range pending {
			batch = append(batch, ev)
		}
		pending = map[string]*protos.MonitorDeviceLocationStreamEventV1_LocationEventV1{}
		return stream.Send(&protos.MonitorDeviceLocationStreamEventV1{Event: batch})
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case c := <-sub.C():
			if len(filter) > 0 {
				if _, ok := filter[c.ev.GetDeviceId()]; !ok {
					continue
				}
			}
			pending[c.ev.GetDeviceId()] = c.ev

		case <-ticker.C:
			if err := flush(); err != nil {
				return err
			}
		}
	}
}
