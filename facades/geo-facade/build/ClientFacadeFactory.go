package build

import (
	bbuild "github.com/Shuv1Wolf/subterra-locate/clients/beacon-admin/build"
	gbuild "github.com/Shuv1Wolf/subterra-locate/clients/geo-renderer/build"
	lbuild "github.com/Shuv1Wolf/subterra-locate/clients/location-engine/build"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
)

type ClientFacadeFactory struct {
	cbuild.CompositeFactory
}

func NewClientFacadeFactory() *ClientFacadeFactory {
	c := &ClientFacadeFactory{
		CompositeFactory: *cbuild.NewCompositeFactory(),
	}

	c.Add(bbuild.NewBeaconsClientFactory())
	c.Add(lbuild.NewLocationEngineClientsFactory())
	c.Add(gbuild.NewGeoRendererClientsFactory())

	return c
}
