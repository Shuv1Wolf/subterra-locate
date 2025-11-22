package containers

import (
	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/build"
	cproc "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
	grpcbuild "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/build"
	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/build"
)

type ZoneProcess struct {
	*cproc.ProcessContainer
}

func NewZoneProcess() *ZoneProcess {
	c := &ZoneProcess{
		ProcessContainer: cproc.NewProcessContainer("zone-processor", "zone-processor microservice"),
	}

	c.AddFactory(build.NewZoneServiceFactory())
	c.AddFactory(grpcbuild.NewDefaultGrpcFactory())
	c.AddFactory(cpg.NewDefaultPostgresFactory())

	return c
}
