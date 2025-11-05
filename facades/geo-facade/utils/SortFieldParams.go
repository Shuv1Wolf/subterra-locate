package utils

import (
	"net/http"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

func GetSortFieldParams(req *http.Request) *cquery.SortField {

	params := req.URL.Query()
	sortParams := make(map[string]string, 0)

	sortParams["name"] = params.Get("name")
	sortParams["ascending"] = params.Get("ascending")

	paging := NewSortFieldFromValue(
		sortParams,
	)
	return paging
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
