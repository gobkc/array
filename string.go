package array

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Has example:
//
//	h := Has(`aaa,bbb,ccc`, `ccc`)
func Has[T string | int | int8 | int32 | int64](parent string, sub T) bool {
	list := strings.Split(parent, ",")
	s := fmt.Sprintf(`%v`, sub)
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func State[CodeOrInt string | int | int8 | int32 | int64, SubType int | string](parent string, sub SubType) CodeOrInt {
	list := strings.Split(parent, ",")
	queryIndex := reflect.TypeOf(sub).Kind() == reflect.String

	s := fmt.Sprintf(`%v`, sub)
	idx, _ := strconv.Atoi(fmt.Sprintf(`%v`, sub))
	for i, v := range list {
		if queryIndex {
			if v == s {
				return CodeOrInt(i)
			}
		} else {
			if i == idx {
				result := *new(CodeOrInt)
				reflect.ValueOf(&result).Elem().Set(reflect.ValueOf(v))
				return result
			}
		}
	}
	return *new(CodeOrInt)
}

func Split(s, sep string, num int) []string {
	container := make([]string, num)
	arr := strings.Split(s, sep)
	for i, str := range arr {
		if i+1 <= num {
			container[i] = str
		}
	}
	return container
}

func QuoteOne(dest string, sign ...string) string {
	var s = `"`
	if len(sign) > 0 {
		s = sign[0]
	}
	return s + dest + s
}
