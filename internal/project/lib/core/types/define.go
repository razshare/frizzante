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

func Define(type_ reflect.Type, ignore []string) (definitions string, root string, known []string, err error) {
	known = ignore
	kind := type_.Kind()
	switch kind {
	case
		reflect.Pointer:
		definitions, root, known, err = Define(type_.Elem(), known)
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

		id := fmt.Sprintf("%s:%s", type_.PkgPath(), type_.Name())

		if slices.Contains(known, id) {
			root = type_.Name()
			return
		}

		if slices.Contains(known, name) {
			parts := strings.Split(type_.PkgPath(), "/")
			count := len(parts)
			name = parts[count-1] + "_" + type_.Name()
		}

		known = append(known, id)
		known = append(known, name)

		var definitionsExtern string

		root = name
		count := type_.NumField()
		definitions += fmt.Sprintf("export type %s = {\n", name)
		for i := 0; i < count; i++ {
			field := type_.Field(i)

			if strings.ToLower(field.Name[0:1]) == field.Name[0:1] {
				continue
			}

			var nameLoc string
			var rootLoc string
			var knownLoc []string
			var definitionsLoc string

			if tag := field.Tag.Get("json"); tag != "" {
				nameLoc = tag
			} else {
				nameLoc = field.Name
			}

			if definitionsLoc, rootLoc, knownLoc, err = Define(field.Type, known); err != nil {
				return
			}

			if rootLoc != "" {
				definitions += fmt.Sprintf("    %s: %s\n", nameLoc, rootLoc)
			}

			if definitionsLoc != "" {
				definitionsExtern += definitionsLoc + "\n\n"
			}

			known = knownLoc
		}
		definitions += "}\n\n"
		definitions += definitionsExtern
		definitions = strings.TrimSpace(definitions)
		root = strings.TrimSpace(root)
	case
		reflect.Slice,
		reflect.Array:
		valueType := type_.Elem()

		var rootLoc string
		var knownLoc []string
		var definitionsLoc string

		if rootLoc, definitionsLoc, knownLoc, err = Define(valueType, known); err != nil {
			return
		}

		root = fmt.Sprintf("%s[]", definitionsLoc)

		if rootLoc != "" {
			definitions = strings.TrimSpace(rootLoc)
		}

		known = knownLoc
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
		var knownLoc []string
		var definitionsLoc string

		if definitionsLoc, rootLoc, knownLoc, err = Define(valueType, known); err != nil {
			return
		}

		root = fmt.Sprintf("Record<%s, %s>", keyTypeName, rootLoc)

		if definitionsLoc != "" {
			definitions = strings.TrimSpace(definitionsLoc)
		}

		known = knownLoc
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
