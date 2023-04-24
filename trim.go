package array

import (
	"reflect"
	"strings"
)

const (
	TrimDefault TrimOption = 0
	WithUpper   TrimOption = 1
	WithLower   TrimOption = 2
	WithUnique  TrimOption = 4
)

type TrimOption = int

func Trim[T any](x T, opts ...TrimOption) T {
	var opt = TrimDefault
	if len(opts) > 0 {
		opt = opts[0]
	}
	toUnique := opt == WithUnique || opt == WithUpper|WithUnique || opt == WithLower|WithUnique
	if reflect.TypeOf(x).Kind() == reflect.Slice && toUnique {
		x = removeDuplicates(x).(T)
	}
	v := reflect.ValueOf(x)
	t := reflect.TypeOf(x)
	nv := reflect.New(t).Elem()
	switch t.Kind() {
	case reflect.String:
		s := strings.TrimSpace(v.String())
		if opt == WithUpper || opt == WithUpper|WithUnique {
			s = strings.ToUpper(s)
		}
		if opt == WithLower || opt == WithLower|WithUnique {
			s = strings.ToLower(s)
		}
		nv = reflect.ValueOf(s)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			nv.Field(i).Set(reflect.ValueOf(Trim(fv.Interface(), opt)))
		}
	case reflect.Slice:
		nv = reflect.MakeSlice(t, v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			ev := v.Index(i)
			nv.Index(i).Set(reflect.ValueOf(Trim(ev.Interface(), opt)))
		}
	default:
		nv = v
	}
	return nv.Interface().(T)
}

func removeDuplicates(slice interface{}) interface{} {
	sliceType := reflect.TypeOf(slice)
	sliceValue := reflect.ValueOf(slice)
	newSliceType := reflect.SliceOf(sliceType.Elem())
	newSliceValue := reflect.MakeSlice(newSliceType, 0, sliceValue.Len())
	seen := make(map[interface{}]bool)
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i).Interface()
		if _, ok := seen[elem]; !ok {
			seen[elem] = true
			newSliceValue = reflect.Append(newSliceValue, reflect.ValueOf(elem))
		}
	}

	return newSliceValue.Interface()
}
