package data1

import (
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

type Map2dV1Schema struct {
	cvalid.ObjectSchema
}

func NewMap2dV1Schema() *Map2dV1Schema {
	c := Map2dV1Schema{}
	c.ObjectSchema = *cvalid.NewObjectSchema()

	c.WithOptionalProperty("id", cconv.String)
	c.WithOptionalProperty("name", cconv.String)
	c.WithOptionalProperty("svg_content", cconv.String)
	c.WithOptionalProperty("scale_x", cconv.Float)
	c.WithOptionalProperty("scale_y", cconv.Float)
	c.WithOptionalProperty("enebled", cconv.Boolean)
	c.WithOptionalProperty("mac_address", cconv.String)
	c.WithOptionalProperty("ip_address", cconv.String)
	return &c
}