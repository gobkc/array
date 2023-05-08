package array

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func MakeSlice[T any](dest ...T) []T {
	objects := make([]T, 0)
	objects = append(objects, dest...)
	return objects
}

var ErrCopySupport = errors.New(`OnlySupportSliceStruct`)

// Copy copies a slice of src to a slice of type []ToType
// It only copies the fields that have the same name or tag
func Copy[ToType any](src any) []ToType {
	typeOf := reflect.TypeOf(src)
	valueOf := reflect.ValueOf(src)
	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	dst := make([]ToType, valueOf.Len(), valueOf.Len())
	for curRow := 0; curRow < valueOf.Len(); curRow++ {
		to := new(ToType)
		s := valueOf.Index(curRow).Interface()
		copyFields(s, to)
		dst[curRow] = *to
	}
	return dst
}

// copyFields copies the fields from src to dst
// It uses reflection to get the field names and tags
func copyFields(src interface{}, dst interface{}) {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)
	srcType := srcVal.Type()
	dstType := dstVal.Type()
	if srcType.Kind() == reflect.Pointer {
		srcType = srcType.Elem()
		srcVal = srcVal.Elem()
	}
	if dstType.Kind() == reflect.Pointer {
		dstType = dstType.Elem()
		dstVal = dstVal.Elem()
	}
	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)
		srcName := srcField.Name
		srcTag := srcField.Tag.Get("json")
		for j := 0; j < dstType.NumField(); j++ {
			dstField := dstType.Field(j)
			dstName := dstField.Name
			dstTag := dstField.Tag.Get("json")
			if dstVal.Field(j).Type().Kind() != srcVal.Field(i).Type().Kind() {
				continue
			}
			if srcName == dstName || (srcTag != `` && dstTag != `` && srcTag == dstTag) {
				dstVal.Field(j).Set(srcVal.Field(i))
				break
			}
		}
	}
}

func Quote[T []string | []int | []int32 | []int64 | float64 | float32](from T, sign ...string) (to []string) {
	var s = `'`
	if len(sign) > 0 {
		s = sign[0]
	}
	valueOf := reflect.ValueOf(from)
	if reflect.TypeOf(from).Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	for curRow := 0; curRow < valueOf.Len(); curRow++ {
		to = append(to, fmt.Sprintf(`%s%v%s`, s, valueOf.Field(curRow).Interface(), s))
	}
	return to
}

func QuoteString[T []string | []int | []int32 | []int64 | float64 | float32](from T, sign ...string) (to string) {
	var s = `'`
	if len(sign) > 0 {
		s = sign[0]
	}
	valueOf := reflect.ValueOf(from)
	if reflect.TypeOf(from).Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	var toSlice []string
	for curRow := 0; curRow < valueOf.Len(); curRow++ {
		toSlice = append(toSlice, fmt.Sprintf(`%s%v%s`, s, valueOf.Field(curRow).Interface(), s))
	}
	return strings.Join(toSlice, `,`)
}

type Element any

func Delete[T Element](slice []T, index int) []T {
	// check if the index is valid
	if index < 0 || index >= len(slice) {
		return slice
	}
	// use the copy function to slide the elements to the left
	copy(slice[index:], slice[index+1:])
	// return the slice without the last element
	return slice[:len(slice)-1]
}

func Remove[T Element](slice []T, element T) []T {
	// loop through the slice to find the element
	for i, v := range slice {
		// if the element is found
		if reflect.DeepEqual(v, element) {
			// use the copy function to slide the elements to the left
			copy(slice[i:], slice[i+1:])
			// return the slice without the last element
			return Remove(slice[:len(slice)-1], element)
		}
	}
	// if the element is not found, return the original slice
	return slice
}

func RemoveFunc[T any](slice []T, match func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if !match(v) {
			result = append(result, v)
		}
	}
	return result
}

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
