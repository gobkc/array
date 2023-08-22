package array

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var hasCache = sync.Map{}

// Has example:
//
//	h := Has(`aaa,bbb,ccc`, `ccc`)
func Has[T string | int | int8 | int32 | int64](parent string, sub T) bool {
	cache, _ := hasCache.Load(parent)
	cacheM, convertOk := cache.(map[any]struct{})
	if convertOk {
		if _, queryOk := cacheM[sub]; queryOk {
			return true
		}
	} else if cacheM == nil {
		cacheM = make(map[any]struct{})
	}
	defer hasCache.Store(parent, cacheM)
	list := strings.Split(parent, ",")
	s := fmt.Sprintf(`%v`, sub)
	for _, v := range list {
		cacheM[v] = struct{}{}
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
