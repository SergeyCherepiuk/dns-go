package utils

import (
	"fmt"
	"reflect"
	"strings"
)

type DiffEntries []DiffEntry

func (e DiffEntries) String() string {
	formattedEntries := make([]string, 0)

	for i, entry := range e {
		formattedEntry := fmt.Sprintf(
			"%d: %s, Actual: %v, Expected: %v",
			i+1,
			entry.FieldPath,
			entry.ActualValue,
			entry.ExpectedValue,
		)
		formattedEntries = append(formattedEntries, formattedEntry)
	}

	return "\n" + strings.Join(formattedEntries, "\n")
}

type DiffEntry struct {
	FieldPath     string
	ActualValue   any
	ExpectedValue any
}

func Diff[T any](actual, expected T) DiffEntries {
	rt := reflect.TypeOf(actual)
	rva := reflect.ValueOf(actual)
	rve := reflect.ValueOf(expected)
	return diff(rt, rva, rve, []string{}, name(rt))
}

func diff(rt reflect.Type, rva, rve reflect.Value, path []string, name string) DiffEntries {
	newPath := copyAppend(path, name)

	switch rt.Kind() {
	case reflect.Func:
		return DiffEntries{}
	case reflect.Slice, reflect.Array:
		return diffSlice(rva, rve, newPath)
	case reflect.Map:
		return diffMap(rva, rve, newPath)
	case reflect.Struct:
		return diffStruct(rt, rva, rve, newPath)
	default:
		if !rva.Equal(rve) {
			return DiffEntries{
				DiffEntry{
					FieldPath:     strings.Join(newPath, "."),
					ActualValue:   rva,
					ExpectedValue: rve,
				},
			}
		}
	}

	return DiffEntries{}
}

func diffSlice(rva, rve reflect.Value, path []string) DiffEntries {
	lenEntries := diffLen(rva, rve, path)
	if len(lenEntries) > 0 {
		return lenEntries
	}

	entries := make(DiffEntries, 0)
	for i := range rva.Len() {
		va := rva.Index(i)
		ve := rve.Index(i)
		fieldEntries := diff(va.Type(), va, ve, path, fmt.Sprintf("[%d]", i))
		entries = append(entries, fieldEntries...)
	}
	return entries
}

func diffMap(rva, rve reflect.Value, path []string) DiffEntries {
	lenEntries := diffLen(rva, rve, path)
	if len(lenEntries) > 0 {
		return lenEntries
	}

	var (
		keys    = rva.MapKeys()
		entries = make(DiffEntries, 0)
	)

	for _, key := range keys {
		va := rva.MapIndex(key)
		ve := rve.MapIndex(key)
		fieldEntries := diff(va.Type(), va, ve, path, fmt.Sprintf("[%v]", key))
		entries = append(entries, fieldEntries...)
	}
	return entries
}

func diffLen(rva, rve reflect.Value, path []string) DiffEntries {
	if rva.Len() != rve.Len() {
		lenPath := copyAppend(path, "len")
		return DiffEntries{
			DiffEntry{
				FieldPath:     strings.Join(lenPath, "."),
				ActualValue:   rva.Len(),
				ExpectedValue: rve.Len(),
			},
		}
	}
	return DiffEntries{}
}

func diffStruct(rt reflect.Type, rva, rve reflect.Value, path []string) DiffEntries {
	entries := make(DiffEntries, 0)
	for i := range rva.Type().NumField() {
		fa := rva.Field(i)
		fe := rve.Field(i)
		fieldEntries := diff(rt.Field(i).Type, fa, fe, path, rt.Field(i).Name)
		entries = append(entries, fieldEntries...)
	}
	return entries
}

func name(rt reflect.Type) string {
	if rt.Name() != "" {
		return rt.Name()
	}

	switch rt.Kind() {
	case reflect.Bool:
		return "<bool>"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "<int>"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "<uint>"
	case reflect.Float32, reflect.Float64:
		return "<float>"
	case reflect.String:
		return "<string>"
	case reflect.Slice:
		return "<slice>"
	case reflect.Array:
		return "<array>"
	case reflect.Map:
		return "<map>"
	case reflect.Struct:
		return "<struct>"
	case reflect.Pointer, reflect.Uintptr, reflect.UnsafePointer:
		return "<pointer>"
	case reflect.Func:
		return "<func>"
	case reflect.Chan:
		return "<chan>"
	default:
		return "<any>"
	}
}

func copyAppend(slice []string, s string) []string {
	sliceCopy := make([]string, len(slice))
	copy(sliceCopy, slice)
	return append(sliceCopy, s)
}
