package utils

import (
	"net/http"
	"reflect"
	"strconv"
	"time"
)

// ParseQueryParams parses query parameters from the request and maps them to a struct
func ParseQueryParams(r *http.Request, out interface{}) error {
	query := r.URL.Query()
	val := reflect.ValueOf(out).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		paramName := fieldType.Tag.Get("query")

		if paramName == "" {
			paramName = fieldType.Name
		}

		paramValue := query.Get(paramName)

		if paramValue == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(paramValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Type() == reflect.TypeOf(time.Time{}) {
				t, err := strconv.ParseInt(paramValue, 10, 64)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(time.Unix(t, 0)))
			} else {
				val, err := strconv.ParseInt(paramValue, 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(val)
			}
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(paramValue, 64)
			if err != nil {
				return err
			}
			field.SetFloat(val)
		default:
			panic("unhandled default case")
		}
	}
	return nil
}
