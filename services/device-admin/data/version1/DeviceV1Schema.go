package data

import (
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

type DeviceV1Schema struct {
	cvalid.ObjectSchema
}

func NewDeviceV1Schema() *DeviceV1Schema {
	c := DeviceV1Schema{}
	c.ObjectSchema = *cvalid.NewObjectSchema()

	c.WithOptionalProperty("id", cconv.String)
	c.WithOptionalProperty("name", cconv.String)
	c.WithOptionalProperty("type", cconv.String)
	c.WithOptionalProperty("model", cconv.String)
	c.WithOptionalProperty("org_id", cconv.String)
	c.WithOptionalProperty("enebled", cconv.Boolean)
	c.WithOptionalProperty("mac_address", cconv.String)
	c.WithOptionalProperty("ip_address", cconv.String)
	return &c
}
