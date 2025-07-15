package containers

import (
	"github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/build"
	cproc "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
	grpcbuild "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/build"
	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/build"
)

type BeaconsProcess struct {
	cproc.ProcessContainer
}

func NewBeaconsProcess() *BeaconsProcess {
	c := &BeaconsProcess{
		ProcessContainer: *cproc.NewProcessContainer("beacon-admin", "beacon-admin microservice"),
	}

	c.AddFactory(build.NewBeaconsServiceFactory())
	c.AddFactory(grpcbuild.NewDefaultGrpcFactory())
	c.AddFactory(cpg.NewDefaultPostgresFactory())

	return c
}
