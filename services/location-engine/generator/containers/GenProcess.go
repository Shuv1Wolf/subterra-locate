package containers

import (
	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/generator/build"
	cproc "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
)

type GenProcess struct {
	*cproc.ProcessContainer
}

func NewGenProcess() *GenProcess {
	c := GenProcess{}
	c.ProcessContainer = cproc.NewProcessContainer("generator", "event gen")

	c.AddFactory(build.NewGenEngineServiceFactory())

	return &c
}
