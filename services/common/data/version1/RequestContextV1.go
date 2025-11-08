package data1

import (
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

type RequestContextV1 struct {
	OrgId  string `json:"org_id"`
	UserId string `json:"user_id"`
}

func NewRequestContextV1Schema() *cvalid.ObjectSchema {
	return cvalid.NewObjectSchema().
		WithOptionalProperty("org_id", cconv.String).
		WithOptionalProperty("user_id", cconv.Boolean)
}

func NewRequestContextV1FromValue(value any) *RequestContextV1 {
	if v, ok := value.(*RequestContextV1); ok {
		return v
	}
	return NewRequestContextV1FromMap(data.NewAnyValueMapFromValue(value))
}

func NewRequestContextV1FromMap(value *data.AnyValueMap) *RequestContextV1 {
	return &RequestContextV1{
		OrgId:  value.GetAsString("org_id"),
		UserId: value.GetAsString("user_id"),
	}
}
