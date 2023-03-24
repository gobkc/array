package array

import (
	"errors"
	"github.com/gobkc/to"
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

var queryIdsMap = map[reflect.Kind]func(dest *reflect.Value, v any){
	reflect.Int: func(dest *reflect.Value, v any) {
		dest.Set(reflect.ValueOf(to.Int[int](v)))
	},
	reflect.Int32: func(dest *reflect.Value, v any) {
		dest.Set(reflect.ValueOf(to.Int[int32](v)))
	},
}

func Ids[T []int | []int32 | []int64 | []float64 | []float32 | []string](pointSlice any, fieldNameOrJsonTag string) *T {
	typeOf := reflect.TypeOf(pointSlice)
	valueOf := reflect.ValueOf(pointSlice)
	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	var result = new(T)
	kind := typeOf.Kind()
	resultTypeOf := reflect.TypeOf(result).Elem()
	newSlice := reflect.MakeSlice(resultTypeOf, valueOf.Len(), valueOf.Len())
	resultKind := newSlice.Index(0).Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return new(T)
	}
	queryFlag := false
	for curRow := 0; curRow < valueOf.Len(); curRow++ {
		rowValueOf := valueOf.Index(curRow)
		rowTypeOf := reflect.TypeOf(rowValueOf.Interface())
		var fName string
		var snakeName string
		for curField := 0; curField < valueOf.Index(curRow).NumField(); curField++ {
			if fName == `` {
				fName = rowTypeOf.Field(curField).Tag.Get("json")
				if fName == "" {
					fName = rowTypeOf.Field(curField).Name
				}
				snakeName = fName
			}
			if fName == fieldNameOrJsonTag || snakeName == fieldNameOrJsonTag {
				queryFlag = true
				v := rowValueOf.Field(curField).Interface()
				switch resultKind {
				case reflect.Int:
					v = to.Any[int](v)
				case reflect.Int32:
					v = to.Any[int32](v)
				case reflect.Int64:
					v = to.Any[int64](v)
				case reflect.Float32:
					v = to.Any[float32](v)
				case reflect.Float64:
					v = to.Any[float64](v)
				case reflect.String:
					v = to.Any[string](v)
				}
				newSlice.Index(curRow).Set(reflect.ValueOf(v))
			}
		}
	}
	if queryFlag {
		reflect.ValueOf(result).Elem().Set(newSlice)
	}
	return result
}
