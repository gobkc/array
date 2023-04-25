package array

import (
	"reflect"
	"strings"
)

func ToSnake(name string) string {
	var convert []byte
	for i, asc := range name {
		if asc >= 65 && asc <= 90 {
			asc += 32
			if i > 0 {
				convert = append(convert, 95)
			}
		}
		convert = append(convert, uint8(asc))
	}
	return string(convert)
}

func CheckIsStructSlice(dest any) bool {
	typeOf := reflect.TypeOf(dest)
	valueOf := reflect.ValueOf(dest)
	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	if typeOf.Kind() != reflect.Slice {
		return false
	}
	tmpSlice := reflect.MakeSlice(typeOf, 1, 1)
	if tmpSlice.Index(0).Kind() == reflect.Pointer {
		pointerKind := tmpSlice.Index(0).Elem().Kind()
		if pointerKind == reflect.Struct {
			return true
		}
		if pointerKind == reflect.Invalid {
			return false
			//unknown := tmpSlice.Index(0).Type().String()
			//switch unknown {
			//case `int`, `int32`, `int64`, `float64`, `float32`, `string`, `*int`, `*int32`, `*int64`, `*float64`, `*float32`, `*string`:
			//	return false
			//case `bool`, `*bool`:
			//	return false
			//default:
			//	return true
			//}
		}
	}
	if tmpSlice.Index(0).Kind() != reflect.Struct {
		return false
	}
	return true
}

func Part(str, sep string, n int) string {
	parts := strings.Split(str, sep)
	n--
	if n < 0 || n >= len(parts) {
		return ""
	}
	return parts[n]
}
