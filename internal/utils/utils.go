package utils

import (
	"fmt"
	"reflect"

	"github.com/bcillie/resy-cli/internal/utils/date"
)

func GetTags(t reflect.Type, tagName string) map[string]string {
	queryTags := make(map[string]string)

	for _, field := range reflect.VisibleFields(t) {
		queryTags[field.Name] = field.Tag.Get(tagName)
	}

	return queryTags
}

func GetQueryParams(t interface{}) map[string]string {
	params := make(map[string]string)

	typeOf := reflect.TypeOf(t)
	valueOf := reflect.ValueOf(t)
	queryTags := GetTags(typeOf, "query")
	for field, prop := range queryTags {
		f := valueOf.FieldByName(field).Interface()
		switch v := f.(type) {
		case int32:
			params[prop] = fmt.Sprint(v)
		case date.ResyDate:
			params[prop] = v.String()
		case string:
			params[prop] = v
		}
	}

	return params
}
