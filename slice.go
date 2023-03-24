package array

import (
	"errors"
	"log"
	"reflect"
)

func MakeSlice[T any](dest ...T) []T {
	objects := make([]T, len(dest))
	objects = append(objects, dest...)
	return objects
}

var ErrCopySupport = errors.New(`OnlySupportSliceStruct`)

func Copy[To any](fromSlice any) *To {
	t := new(To)
	typeOf := reflect.TypeOf(fromSlice)
	valueOf := reflect.ValueOf(fromSlice)
	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	if !CheckIsStructSlice(t) || !CheckIsStructSlice(valueOf.Interface()) {
		log.Println(`Array.Copy:`, ErrCopySupport)
		return nil
	}
	to := reflect.MakeSlice(reflect.TypeOf(t).Elem(), 1, 1)
	newTo := to.Index(0)
	toMap := make(map[string]int)
	toTypeOf := newTo.Type()
	for i := 0; i < reflect.ValueOf(newTo).NumField(); i++ {
		fName := toTypeOf.Field(i).Name
		toMap[fName] = i
		tag := toTypeOf.Field(i).Tag.Get(`json`)
		if tag != `` {
			toMap[tag] = i
		}
		toMap[ToSnake(fName)] = i
	}
	var values = reflect.MakeSlice(reflect.TypeOf(t).Elem(), valueOf.Len(), valueOf.Len())
	for curRow := 0; curRow < valueOf.Len(); curRow++ {
		rowValueOf := valueOf.Index(curRow)
		rowTypeOf := reflect.TypeOf(rowValueOf.Interface())
		toRow := newTo
		for curField := 0; curField < valueOf.Index(curRow).NumField(); curField++ {
			fName := rowTypeOf.Field(curField).Tag.Get("json")
			if fName == "" {
				fName = rowTypeOf.Field(curField).Name
			}
			toField, ok := toMap[fName]
			if !ok {
				continue
			}
			if toName := toTypeOf.Field(toField).Name; toName == fName || ToSnake(toName) == fName ||
				toTypeOf.Field(toField).Tag.Get(`json`) == fName {
				if toRow.Field(toField).Kind() == rowValueOf.Field(curField).Kind() &&
					toRow.Field(toField).Kind() != reflect.Pointer && rowValueOf.Field(curField).Kind() != reflect.Pointer {
					toRow.Field(toField).Set(rowValueOf.Field(curField))
				}
			}
		}
		values.Index(curRow).Set(toRow)
	}
	reflect.ValueOf(t).Elem().Set(values)
	return t
}
