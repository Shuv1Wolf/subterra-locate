package clients1

import (
	"context"

	protos "github.com/Shuv1Wolf/subterra-locate/services/zone-processor/protos"
	cclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
	"google.golang.org/grpc"
)

type ZoneMonitorGrpcClientV1 struct {
	*cclients.GrpcClient
	client protos.ZoneMonitorClient
}

func NewZoneMonitorGrpcClientV1() *ZoneMonitorGrpcClientV1 {
	c := &ZoneMonitorGrpcClientV1{
		GrpcClient: cclients.NewGrpcClient("zone.monitor.v1"),
	}
	return c
}

func (c *ZoneMonitorGrpcClientV1) Open(ctx context.Context) error {
	err := c.GrpcClient.Open(ctx)
	if err == nil {
		c.client = protos.NewZoneMonitorClient(c.Connection)
	}
	return err
}

func (c *ZoneMonitorGrpcClientV1) MonitorZoneV1(ctx context.Context, orgId string, mapId string) (grpc.ServerStreamingClient[protos.MonitorZoneStreamEventV1], error) {
	request := &protos.MonitorZoneRequestV1{
		OrgId: orgId,
		MapId: mapId,
	}

	stream, err := c.client.MonitorZoneV1(ctx, request)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
