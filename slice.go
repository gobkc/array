package array

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Slices list := array.Slices(item)
func Slices[T any](dest ...T) []T {
	objects := make([]T, 0)
	for _, v := range dest {
		if !reflect.ValueOf(v).IsZero() {
			objects = append(objects, v)
		}
	}
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

func SliceToPointer[T any](s []T) []*T {
	p := make([]*T, len(s), cap(s))
	for i := range s {
		p[i] = &s[i]
	}
	return p
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

// Quote Example
// Quote([]string{"aa", "bbb"}) result:['aa' 'bbb']
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
		to = append(to, fmt.Sprintf(`%s%v%s`, s, valueOf.Index(curRow).Interface(), s))
	}
	return to
}

// QuoteString Example
// QuoteString([]string{"aa", "bbb"}) result:'aa','bbb'
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
		toSlice = append(toSlice, fmt.Sprintf(`%s%v%s`, s, valueOf.Index(curRow).Interface(), s))
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

type KeyTypeDef interface {
	string | int8 | int | int32 | int64 | float32 | float64
}

// MakeMap slice convert map
// appMaps := array.MakeMap[int64](apps, `plan_id`)
// appMaps result: map[plan_id]app
func MakeMap[KeyType KeyTypeDef, ItemType any](slice []ItemType, field string) map[KeyType]ItemType {
	result := make(map[KeyType]ItemType)
	for _, elem := range slice {
		v := reflect.ValueOf(elem)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Struct {
			t := v.Type()
			index := -1
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if f.Name == field || f.Tag.Get("json") == field {
					index = i
					break
				}
			}
			if index >= 0 {
				key := v.Field(index).Interface()
				if convertKey, ok := key.(KeyType); ok {
					result[convertKey] = elem
				}
			}
		}
	}
	return result
}

// MakeMaps slice convert map
// Example :
// appMaps := array.MakeMaps[int64](apps, `plan_id`)
// appMaps result: map[plan_id][]apps
func MakeMaps[KeyType KeyTypeDef, ItemType any](slice []ItemType, field string) map[KeyType][]ItemType {
	result := make(map[KeyType][]ItemType)
	for _, elem := range slice {
		v := reflect.ValueOf(elem)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Struct {
			t := v.Type()
			index := -1
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if f.Name == field || f.Tag.Get("json") == field {
					index = i
					break
				}
			}
			if index >= 0 {
				key := v.Field(index).Interface()
				if convertKey, ok := key.(KeyType); ok {
					if query, queryOk := result[convertKey]; queryOk {
						result[convertKey] = append(query, elem)
					} else {
						result[convertKey] = []ItemType{elem}
					}
				}
			}
		}
	}
	return result
}

// Index get slice elem
// Example:
//
//	 a := []func(){
//			func() {
//				fmt.Println(123)
//			},
//			func() {
//				fmt.Println(456)
//			},
//		}
//		ss := Index(a, 2)
//		if ss != nil {
//			ss()
//		} else {
//			fmt.Println(`ss is nil`)
//		}
func Index[T any](list []T, index int) T {
	var result = new(T)
	if len(list) <= index || index < 0 {
		return *result
	}
	return list[index]
}

// Count example:
//
//	count := Count(ss, func(i int, row *T2) bool {
//			return row.Name == `A1`
//	})
func Count[T any](array []T, fc func(i int, row T) bool) int {
	count := 0
	for i, row := range array {
		if fc(i, row) {
			count++
		}
	}
	return count
}

// Pack example
//
//	d := Pack(si, func(i int, item *SI) bool {
//		return item.Age > 1
//	})
func Pack[T any](slices []T, f func(i int, item T) bool) (dest []T) {
	for i, slice := range slices {
		if ok := f(i, slice); ok {
			dest = append(dest, slice)
		}
	}
	return
}

// PackField example:
//
//	d := PackField[string](si, func(i int, item *SI) string {
//		if item.Age > 1 {
//			return item.Name
//		}
//		return ``
//	})
func PackField[F, T any](slices []T, f func(i int, item T) F) (dest []F) {
	for i, slice := range slices {
		field := f(i, slice)
		itemType := reflect.TypeOf(field)
		itemValue := reflect.ValueOf(field)
		if itemType.Kind() == reflect.Pointer {
			itemValue = itemValue.Elem()
		}
		isZero := itemValue.IsZero()

		if !isZero {
			dest = append(dest, field)
		}
	}
	return
}

// QueryElement example:
// fmt.Println(QueryElement[string](si, "0.Schools.Middle.Name.1"))
func QueryElement[T any](dest any, path string) T {
	parts := strings.Split(strings.TrimPrefix(path, "."), ".")
	// Get the reflect.Value of dest
	v := reflect.ValueOf(dest)
	// Loop through the parts of the path
	for _, part := range parts {
		// If v is invalid or nil, return the zero value of T
		if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
			return reflect.Zero(reflect.TypeOf((*T)(nil)).Elem()).Interface().(T)
		}
		// If v is a pointer, dereference it
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		// Switch on the kind of v
		switch v.Kind() {
		case reflect.Struct:
			// If v is a struct, use FieldByName to get the field with the name of part
			if _, ok := v.Type().FieldByName(part); ok {
				v = v.FieldByName(part)
			} else {
				return reflect.Zero(reflect.TypeOf((*T)(nil)).Elem()).Interface().(T)
			}
		case reflect.Slice, reflect.Array:
			index, _ := strconv.Atoi(part)
			if v.Len() > index {
				v = v.Index(index)
			} else {
				return reflect.Zero(reflect.TypeOf((*T)(nil)).Elem()).Interface().(T)
			}
		case reflect.Map:
			// If v is a map, use MapIndex to get the value with the key of part
			v = v.MapIndex(reflect.ValueOf(part))
		default:
			// If v is none of the above, return the zero value of T
			return reflect.Zero(reflect.TypeOf((*T)(nil)).Elem()).Interface().(T)
		}
	}
	// Return the final value of v as T
	return v.Interface().(T)
}
