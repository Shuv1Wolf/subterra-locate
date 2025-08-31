package version1

import (
	"context"

	protos "github.com/Shuv1Wolf/subterra-locate/services/location-engine/protos"
	cclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
	"google.golang.org/grpc"
)

type LocationMonitorGrpcClientV1 struct {
	*cclients.GrpcClient
	client protos.LocationMonitorClient
}

func NewLocationMonitorGrpcClientV1() *LocationMonitorGrpcClientV1 {
	c := &LocationMonitorGrpcClientV1{
		GrpcClient: cclients.NewGrpcClient("location.monitor.v1"),
	}
	return c
}

func (c *LocationMonitorGrpcClientV1) Open(ctx context.Context, correlationId string) error {
	err := c.GrpcClient.Open(ctx)
	if err == nil {
		c.client = protos.NewLocationMonitorClient(c.Connection)
	}
	return err
}

func (c *LocationMonitorGrpcClientV1) MonitorDeviceLocation(ctx context.Context, orgId string, deviceIds []string) (grpc.ServerStreamingClient[protos.MonitorDeviceLocationStreamEventV1], error) {
	request := &protos.MonitorDeviceLocationRequestV1{
		OrgId:    orgId,
		DeviceId: deviceIds,
	}

	stream, err := c.client.MonitorDeviceLocationV1(ctx, request)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
