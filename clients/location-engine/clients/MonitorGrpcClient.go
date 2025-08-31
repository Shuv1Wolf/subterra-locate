package clients

import (
	protos "github.com/Shuv1Wolf/subterra-locate/services/location-engine/protos"
)

type MonitorGrpcClient struct {
	client protos.MonitorClient
}
