package data1

import (
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

type ZoneV1Schema struct {
	cvalid.ObjectSchema
}

func NewZoneV1Schema() *ZoneV1Schema {
	c := ZoneV1Schema{}
	c.ObjectSchema = *cvalid.NewObjectSchema()

	c.WithOptionalProperty("id", cconv.String)
	c.WithOptionalProperty("map_id", cconv.String)
	c.WithOptionalProperty("org_id", cconv.String)
	c.WithOptionalProperty("name", cconv.String)
	c.WithOptionalProperty("position_x", cconv.Float)
	c.WithOptionalProperty("position_y", cconv.Float)
	c.WithOptionalProperty("width", cconv.Float)
	c.WithOptionalProperty("height", cconv.Float)
	c.WithOptionalProperty("type", cconv.String)
	return &c
}
