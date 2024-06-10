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

// TODO: Unit test
func Diff[T any](actual, expected T) DiffEntries {
	var (
		rva  = reflect.ValueOf(actual)
		rve  = reflect.ValueOf(expected)
		path = []string{rva.Type().Name()}
	)

	return diff(rva, rve, path)
}

func diff(rva, rve reflect.Value, path []string) []DiffEntry {
	var (
		rt      = rva.Type()
		entries = make([]DiffEntry, 0)
	)

	for i := range rt.NumField() {
		var (
			fa = rva.Field(i)
			fe = rve.Field(i)
		)

		if rt.Field(i).Type.Kind() == reflect.Struct {
			substructPath := copyAppend(path, rt.Field(i).Name)
			substructEntries := diff(fa, fe, substructPath)
			entries = append(entries, substructEntries...)
			continue
		}

		if !fa.Equal(fe) {
			fieldPath := copyAppend(path, rt.Field(i).Name)
			entry := DiffEntry{
				FieldPath:     strings.Join(fieldPath, "."),
				ActualValue:   fa,
				ExpectedValue: fe,
			}
			entries = append(entries, entry)
		}
	}

	return entries
}

func copyAppend(slice []string, s string) []string {
	sliceCopy := make([]string, len(slice))
	copy(sliceCopy, slice)
	return append(sliceCopy, s)
}
