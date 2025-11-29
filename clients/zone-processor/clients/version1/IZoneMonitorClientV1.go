package clients1

import (
	"context"

	protos "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/protos"
	"google.golang.org/grpc"
)

type IZoneMonitorClientV1 interface {
	MonitorZoneV1(ctx context.Context, orgId string, mapId string) (grpc.ServerStreamingClient[protos.MonitorZoneStreamEventV1], error)
}
