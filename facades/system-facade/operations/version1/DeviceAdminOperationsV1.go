package operations1

import (
	"context"
	"net/http"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	httpcontr "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"

	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/device-admin/clients/version1"
	services1 "github.com/Shuv1Wolf/subterra-locate/services/device-admin/data/version1"
)

type DeviceAdminOperationsV1 struct {
	*httpcontr.RestOperations
	deviceAdmin clients1.IDeviceClientV1
}

func NewDeviceAdminOperationsV1() *DeviceAdminOperationsV1 {
	c := DeviceAdminOperationsV1{
		RestOperations: httpcontr.NewRestOperations(),
	}
	c.DependencyResolver.Put(context.Background(), "device-admin", cref.NewDescriptor("device-admin", "client", "*", "*", "1.0"))
	return &c
}

func (c *DeviceAdminOperationsV1) SetReferences(ctx context.Context, references cref.IReferences) {
	c.RestOperations.SetReferences(ctx, references)

	dependency, _ := c.DependencyResolver.GetOneRequired("device-admin")
	client, ok := dependency.(clients1.IDeviceClientV1)
	if !ok {
		panic("DeviceAdminOperationsV1: Cant't resolv dependency 'client' to IDeviceClientV1")
	}
	c.deviceAdmin = client
}

func (c *DeviceAdminOperationsV1) GetDevices(res http.ResponseWriter, req *http.Request) {
	var filter = c.GetFilterParams(req)
	var paging = c.GetPagingParams(req)

	page, err := c.deviceAdmin.GetDevices(
		context.Background(), *filter, *paging)

	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, page, nil)
	}
}

func (c *DeviceAdminOperationsV1) GetDeviceById(res http.ResponseWriter, req *http.Request) {
	id := c.GetParam(req, "id")
	item, err := c.deviceAdmin.GetDeviceById(context.Background(), id)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *DeviceAdminOperationsV1) CreateDevice(res http.ResponseWriter, req *http.Request) {

	data := services1.DeviceV1{}
	err := c.DecodeBody(req, &data)
	if err != nil {
		c.SendError(res, req, err)
	}
	item, err := c.deviceAdmin.CreateDevice(context.Background(), data)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *DeviceAdminOperationsV1) UpdateDevice(res http.ResponseWriter, req *http.Request) {
	data := services1.DeviceV1{}
	err := c.DecodeBody(req, &data)
	if err != nil {
		c.SendError(res, req, err)
	}

	item, err := c.deviceAdmin.UpdateDevice(context.Background(), data)
	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}

func (c *DeviceAdminOperationsV1) DeleteDeviceById(res http.ResponseWriter, req *http.Request) {
	id := c.GetParam(req, "id")

	item, err := c.deviceAdmin.DeleteDeviceById(context.Background(), id)

	if err != nil {
		c.SendError(res, req, err)
	} else {
		c.SendResult(res, req, item, nil)
	}
}
