package service

import (
	"time"

	protos "github.com/Shuv1Wolf/subterra-locate/services/location-engine/protos"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/utils"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	grpc "google.golang.org/grpc"
)

const heartbeatInterval = 5 * time.Second

type MonitorLocation struct {
	protos.UnimplementedLocationMonitorServer
	deviceState      *utils.DeviceStateStore
	beaconStateStore *utils.BeaconStateStore
	logger           *clog.CompositeLogger
}

func NewMonitorLocation(deviceState *utils.DeviceStateStore, beaconStateStore *utils.BeaconStateStore, logger *clog.CompositeLogger) *MonitorLocation {
	return &MonitorLocation{
		deviceState:      deviceState,
		beaconStateStore: beaconStateStore,
		logger:           logger,
	}
}

func (s *MonitorLocation) MonitorDeviceLocationV1(
	in *protos.MonitorDeviceLocationRequestV1,
	stream grpc.ServerStreamingServer[protos.MonitorDeviceLocationStreamEventV1],
) error {
	ctx := stream.Context()
	orgID := in.GetOrgId()
	mapID := in.GetMapId()

	devFilter := make(map[string]struct{}, len(in.GetDeviceId()))
	for _, id := range in.GetDeviceId() {
		devFilter[id] = struct{}{}
	}

	initial := s.deviceState.Snapshot(orgID)
	resp := make([]*protos.MonitorDeviceLocationStreamEventV1_LocationEventV1, 0, len(initial))
	for _, ev := range initial {
		if mapID != "" && ev.GetMapId() != mapID {
			continue
		}
		if len(devFilter) > 0 {
			if _, ok := devFilter[ev.GetDeviceId()]; !ok {
				continue
			}
		}
		resp = append(resp, ev)
	}
	if err := stream.Send(&protos.MonitorDeviceLocationStreamEventV1{Event: resp}); err != nil {
		return err
	}

	sub := s.deviceState.Subscribe(orgID)
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
			if len(devFilter) > 0 {
				if _, ok := devFilter[c.Ev.GetDeviceId()]; !ok {
					continue
				}
			}

			if c.Ev.X == 0 && c.Ev.Y == 0 && c.Ev.Z == 0 {
				pending[c.Ev.GetDeviceId()] = c.Ev
			}

			if mapID != "" && c.Ev.GetMapId() != mapID {
				continue
			}

			pending[c.Ev.GetDeviceId()] = c.Ev

		case <-ticker.C:
			if err := flush(); err != nil {
				return err
			}
		}
	}
}

func (s *MonitorLocation) MonitorBeaconLocationV1(
	in *protos.MonitorBeaconLocationRequestV1,
	stream grpc.ServerStreamingServer[protos.MonitorBeaconLocationStreamEventV1],
) error {
	ctx := stream.Context()
	orgID := in.GetOrgId()
	mapID := in.GetMapId()

	beaconFilter := make(map[string]struct{}, len(in.GetBeaconId()))
	for _, id := range in.GetBeaconId() {
		beaconFilter[id] = struct{}{}
	}

	initial := s.beaconStateStore.Snapshot(orgID)
	resp := make([]*protos.MonitorBeaconLocationStreamEventV1_LocationEventV1, 0, len(initial))
	for _, ev := range initial {
		if mapID != "" && ev.GetMapId() != mapID {
			continue
		}
		if len(beaconFilter) > 0 {
			if _, ok := beaconFilter[ev.GetBeaconId()]; !ok {
				continue
			}
		}
		resp = append(resp, ev)
	}
	if err := stream.Send(&protos.MonitorBeaconLocationStreamEventV1{Event: resp}); err != nil {
		return err
	}

	sub := s.beaconStateStore.Subscribe(orgID)
	defer sub.Close()

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	pending := map[string]*protos.MonitorBeaconLocationStreamEventV1_LocationEventV1{}

	flush := func() error {
		if len(pending) == 0 {
			return stream.Send(&protos.MonitorBeaconLocationStreamEventV1{
				Event: nil,
			})
		}
		batch := make([]*protos.MonitorBeaconLocationStreamEventV1_LocationEventV1, 0, len(pending))
		for _, ev := range pending {
			batch = append(batch, ev)
		}
		pending = map[string]*protos.MonitorBeaconLocationStreamEventV1_LocationEventV1{}
		return stream.Send(&protos.MonitorBeaconLocationStreamEventV1{Event: batch})
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case c := <-sub.C():
			if len(beaconFilter) > 0 {
				if _, ok := beaconFilter[c.Ev.GetBeaconId()]; !ok {
					continue
				}
			}

			if c.Ev.X == 0 && c.Ev.Y == 0 && c.Ev.Z == 0 {
				pending[c.Ev.GetBeaconId()] = c.Ev
			}

			if mapID != "" && c.Ev.GetMapId() != mapID {
				continue
			}

			pending[c.Ev.GetBeaconId()] = c.Ev

		case <-ticker.C:
			if err := flush(); err != nil {
				return err
			}
		}
	}
}
