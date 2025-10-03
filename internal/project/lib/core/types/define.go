package types

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func IsPrimitive(type_ reflect.Type) bool {
	switch type_.Kind() {
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex128,
		reflect.Uintptr,
		reflect.String:
		return true
	default:
		return false
	}
}

func Define(type_ reflect.Type, packages map[string][]string, definitions map[string]map[string][]string) (root string, err error) {
	kind := type_.Kind()
	switch kind {
	case
		reflect.Pointer:
		root, err = Define(type_.Elem(), packages, definitions)
		return
	case
		reflect.Struct,
		reflect.Interface:
		name := type_.Name()

		if name == "" {
			root = "unknown"
			return
		}

		if name == "error" {
			root = "string"
			return
		}

		parts := strings.Split(type_.PkgPath(), "/")
		count := len(parts)
		package_ := parts[count-1]
		if _, exists := definitions[package_]; !exists {
			definitions[package_] = map[string][]string{}
		}

		if _, exists := definitions[package_][name]; !exists {
			definitions[package_][name] = make([]string, 0)
		}

		if _, exists := packages[package_]; exists {
			if slices.Contains(packages[package_], name) {
				root = package_ + "." + name
				return
			}
			packages[package_] = append(packages[package_], name)
		} else {
			packages[package_] = []string{name}
		}

		root = package_ + "." + name
		count = type_.NumField()
		for i := 0; i < count; i++ {
			field := type_.Field(i)

			if strings.ToLower(field.Name[0:1]) == field.Name[0:1] {
				continue
			}

			var nameLoc string
			var rootLoc string

			if tag := field.Tag.Get("json"); tag != "" {
				nameLoc = tag
			} else {
				nameLoc = field.Name
			}

			if rootLoc, err = Define(field.Type, packages, definitions); err != nil {
				return
			}

			if rootLoc != "" {
				definitions[package_][name] = append(definitions[package_][name], fmt.Sprintf("%s: %s", nameLoc, rootLoc))
			}
		}
		root = strings.TrimSpace(root)
	case
		reflect.Slice,
		reflect.Array:
		valueType := type_.Elem()

		var rootLoc string

		if rootLoc, err = Define(valueType, packages, definitions); err != nil {
			return
		}

		root = fmt.Sprintf("%s[]", rootLoc)
	case
		reflect.Map:
		keyType := type_.Key()
		keyTypeName := keyType.Name()

		if !IsPrimitive(keyType) {
			err = errors.New("map key type must be primitive")
			return
		}

		valueType := type_.Elem()

		var rootLoc string

		if rootLoc, err = Define(valueType, packages, definitions); err != nil {
			return
		}

		root = fmt.Sprintf("Record<%s, %s>", keyTypeName, rootLoc)
	case
		reflect.Chan,
		reflect.Func,
		reflect.Invalid,
		reflect.UnsafePointer:
		err = fmt.Errorf("type %s of kind %s is not supported", type_.String(), kind.String())
	case
		reflect.Bool:
		root = "boolean"
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex128,
		reflect.Uintptr:
		root = "number"
	case
		reflect.String:
		root = "string"
	default:
		root = "unknown"
	}
	return
}
