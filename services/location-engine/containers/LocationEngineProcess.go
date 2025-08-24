package containers

import (
	bbuild "github.com/Shuv1Wolf/subterra-locate/clients/beacon-admin/build"
	dbuild "github.com/Shuv1Wolf/subterra-locate/clients/device-admin/build"
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/build"
	cproc "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
)

type LocationEngineProcess struct {
	*cproc.ProcessContainer
}

func NewLocationEngineProcess() *LocationEngineProcess {
	c := LocationEngineProcess{}
	c.ProcessContainer = cproc.NewProcessContainer("location-engine", "Location Engine")

	c.AddFactory(build.NewLocationEngineServiceFactory())
	c.AddFactory(bbuild.NewBeaconsClientFactory())
	c.AddFactory(dbuild.NewDeviceClientFactory())

	return &c
}
