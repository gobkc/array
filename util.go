package array

import "reflect"

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
	if tmpSlice.Index(0).Kind() != reflect.Struct {
		return false
	}
	return true
}
