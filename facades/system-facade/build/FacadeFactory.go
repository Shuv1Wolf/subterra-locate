package build

import (
	controllers1 "github.com/Shuv1Wolf/subterra-locate/facades/system-facade/controllers/version1"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type FacadeFactory struct {
	cbuild.Factory
	FacadeServiceV1Descriptor *cref.Descriptor
}

func NewFacadeFactory() *FacadeFactory {

	c := FacadeFactory{
		Factory: *cbuild.NewFactory(),
	}
	c.FacadeServiceV1Descriptor = cref.NewDescriptor("system-facade", "controller", "http", "*", "1.0")
	c.RegisterType(c.FacadeServiceV1Descriptor, controllers1.NewFacadeControllerV1)
	return &c
}
