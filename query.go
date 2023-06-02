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

func First[T any](s []T) T {
	if len(s) == 0 {
		t := new(T)
		return *t
	}
	return s[0]
}

var queryIdsMap = map[reflect.Kind]func(dest *reflect.Value, v any){
	reflect.Int: func(dest *reflect.Value, v any) {
		dest.Set(reflect.ValueOf(to.Int[int](v)))
	},
	reflect.Int32: func(dest *reflect.Value, v any) {
		dest.Set(reflect.ValueOf(to.Int[int32](v)))
	},
}

func Ids[T int | int32 | int64 | float64 | float32 | string](pointSlice any, fieldNameOrJsonTag string) *[]T {
	result := make([]T, 0)
	psType := reflect.TypeOf(pointSlice)
	psValue := reflect.ValueOf(pointSlice)
	if psType.Kind() != reflect.Slice {
		return &result
	}
	for i := 0; i < psValue.Len(); i++ {
		elem := psValue.Index(i)
		elemType := elem.Type()
		if elemType.Kind() == reflect.Pointer {
			elem = elem.Elem()
			elemType = elemType.Elem()
		}
		if elemType.Kind() != reflect.Struct {
			continue
		}
		var field reflect.Value
		if _, ok := elemType.FieldByName(fieldNameOrJsonTag); ok {
			field = elem.FieldByName(fieldNameOrJsonTag)
		} else {
			for j := 0; j < elemType.NumField(); j++ {
				f := elemType.Field(j)
				if tag, ok := f.Tag.Lookup("json"); ok && tag == fieldNameOrJsonTag {
					field = elem.Field(j)
					break
				}
			}
		}
		if field.IsValid() {
			var v any
			var t any = (*T)(nil)
			switch t.(type) {
			case *int:
				v = to.Any[int](field.Interface())
			case *int32:
				v = to.Any[int32](field.Interface())
			case *int64:
				v = to.Any[int64](field.Interface())
			case *string:
				v = to.Any[string](field.Interface())
			case *float32:
				v = to.Any[float32](field.Interface())
			case *float64:
				v = to.Any[float64](field.Interface())
			default:
				continue
			}
			result = append(result, v.(T))
		}
	}
	return &result
}
