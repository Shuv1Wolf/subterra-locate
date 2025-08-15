package data

import (
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

type BeaconV1Schema struct {
	cvalid.ObjectSchema
}

func NewBeaconV1Schema() *BeaconV1Schema {
	c := BeaconV1Schema{}
	c.ObjectSchema = *cvalid.NewObjectSchema()

	c.WithOptionalProperty("id", cconv.String)
	c.WithOptionalProperty("type", cconv.String)
	c.WithRequiredProperty("udi", cconv.String)
	c.WithOptionalProperty("label", cconv.String)
	c.WithOptionalProperty("x", cconv.Float)
	c.WithOptionalProperty("y", cconv.Float)
	c.WithOptionalProperty("z", cconv.Float)
	c.WithOptionalProperty("org_id", cconv.String)
	c.WithOptionalProperty("enabled", cconv.Boolean)
	return &c
}
