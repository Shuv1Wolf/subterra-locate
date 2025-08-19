package containers

import (
	"github.com/Shuv1Wolf/subterra-locate/services/device-admin/build"
	cproc "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
	grpcbuild "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/build"
	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/build"
)

type DeviceProcess struct {
	cproc.ProcessContainer
}

func NewDeviceProcess() *DeviceProcess {
	c := &DeviceProcess{
		ProcessContainer: *cproc.NewProcessContainer("device-admin", "device-admin microservice"),
	}

	c.AddFactory(build.NewDeviceServiceFactory())
	c.AddFactory(grpcbuild.NewDefaultGrpcFactory())
	c.AddFactory(cpg.NewDefaultPostgresFactory())

	return c
}
