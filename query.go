package array

import (
	"errors"
	"reflect"
)

func Foreach[T comparable](list []T, f func(line int, row *T)) {
	for i, item := range list {
		f(i, &item)
		list[i] = item
	}
}

var ErrEmptyArray = errors.New("ErrEmptyArray")

func First[T any](array []T) (*T, error) {
	if len(array) == 0 {
		return nil, ErrEmptyArray
	}
	first := &array[0]
	return first, nil
}

func Ids[T []int | []int32 | []int64 | []float64 | []float32 | []string](pointSlice any, fieldNameOrJsonTag string) *T {
	var result = new(T)
	resultTypeOf := reflect.TypeOf(result).Elem()
	newSlice := reflect.MakeSlice(resultTypeOf, 1, 1)

	typeOf := reflect.TypeOf(pointSlice)
	valueOf := reflect.ValueOf(pointSlice)
	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	kind := typeOf.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return result
	}
	var fieldKind reflect.Kind
	var values []reflect.Value
	for curRow := 0; curRow < valueOf.Len(); curRow++ {
		rowValueOf := valueOf.Index(curRow)
		rowTypeOf := reflect.TypeOf(rowValueOf.Interface())
		for curField := 0; curField < valueOf.Index(curRow).NumField(); curField++ {
			fName := rowTypeOf.Field(curField).Tag.Get("json")
			if fName == "" {
				fName = rowTypeOf.Field(curField).Name
			}
			if fieldKind == reflect.Invalid {
				fieldKind = rowValueOf.Field(curField).Kind()
			}
			if fName == fieldNameOrJsonTag {
				v := rowValueOf.Field(curField).Interface()
				switch fieldKind {
				case reflect.Int32:
					values = append(values, reflect.ValueOf(to.Int[int32](v)))
				}
			}
		}
	}
	return result
}
