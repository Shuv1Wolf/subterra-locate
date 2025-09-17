package containers

import (
	"github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/build"
	cproc "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
	grpcbuild "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/build"
	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/build"
)

type MapProcess struct {
	cproc.ProcessContainer
}

func NewMapProcess() *MapProcess {
	c := &MapProcess{
		ProcessContainer: *cproc.NewProcessContainer("geo-renderer", "geo-renderer microservice"),
	}

	c.AddFactory(build.NewMapServiceFactory())
	c.AddFactory(grpcbuild.NewDefaultGrpcFactory())
	c.AddFactory(cpg.NewDefaultPostgresFactory())

	return c
}
