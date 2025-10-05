package clients1

import (
	"context"

	protos "github.com/Shuv1Wolf/subterra-locate/services/location-engine/protos"
	"google.golang.org/grpc"
)

type ILocationMonitorClientV1 interface {
	MonitorDeviceLocation(ctx context.Context, orgId string, mapId string, deviceIds []string) (grpc.ServerStreamingClient[protos.MonitorDeviceLocationStreamEventV1], error)
	MonitorBeaconLocation(ctx context.Context, orgId string, mapId string, beaconIds []string) (grpc.ServerStreamingClient[protos.MonitorBeaconLocationStreamEventV1], error)
}
