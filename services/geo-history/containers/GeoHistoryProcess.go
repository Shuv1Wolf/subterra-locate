package containers

import (
	"github.com/Shuv1Wolf/subterra-locate/services/geo-history/build"
	cproc "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
	grpcbuild "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/build"
	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/build"
)

type GeoHistoryProcess struct {
	cproc.ProcessContainer
}

func NewGeoHistoryProcess() *GeoHistoryProcess {
	c := &GeoHistoryProcess{
		ProcessContainer: *cproc.NewProcessContainer("geo-history", "geo-history microservice"),
	}

	c.AddFactory(build.NewGeoHistoryServiceFactory())
	c.AddFactory(grpcbuild.NewDefaultGrpcFactory())
	c.AddFactory(cpg.NewDefaultPostgresFactory())

	return c
}
