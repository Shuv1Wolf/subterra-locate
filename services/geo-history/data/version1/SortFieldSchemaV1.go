package data1

import (
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

func NewSortFieldSchema() *cvalid.ObjectSchema {
	return cvalid.NewObjectSchema().
		WithOptionalProperty("name", cconv.String).
		WithOptionalProperty("ascending", cconv.Boolean)
}

func NewSortFieldFromValue(value any) *cquery.SortField {
	if v, ok := value.(*cquery.SortField); ok {
		return v
	}
	return NewSortFieldFromMap(data.NewAnyValueMapFromValue(value))
}

func NewSortFieldFromMap(value *data.AnyValueMap) *cquery.SortField {
	return &cquery.SortField{
		Name:      value.GetAsStringWithDefault("name", "timestamp"),
		Ascending: value.GetAsBooleanWithDefault("ascending", true),
	}
}
