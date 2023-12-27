package utils

import (
	"fmt"
	"reflect"
	"time"
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
	fmtTags := GetTags(typeOf, "fmt")
	for field, prop := range queryTags {
		f := valueOf.FieldByName(field).Interface()
		switch v := f.(type) {
		case int32:
			params[prop] = fmt.Sprint(v)
		case time.Time:
			params[prop] = v.Format(fmtTags[field])
		case string:
			params[prop] = v
		}
		fmt.Printf("%s: %s\n", prop, params[prop])
	}

	return params
}
