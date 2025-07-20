package container

import (
	"github.com/Shuv1Wolf/subterra-locate/facades/geo-facade/build"
	cproc "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
	httpbuild "github.com/pip-services4/pip-services4-go/pip-services4-http-go/build"
)

type FacadeProcess struct {
	*cproc.ProcessContainer
}

func NewFacadeProcess() *FacadeProcess {

	c := FacadeProcess{}
	c.ProcessContainer = cproc.NewProcessContainer("geo-facade", "Public facade for geo services")
	c.AddFactory(build.NewClientFacadeFactory())
	c.AddFactory(build.NewFacadeFactory())
	c.AddFactory(httpbuild.NewDefaultHttpFactory())

	return &c
}
