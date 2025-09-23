package types

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func Extract(_type reflect.Type, ignore []string) (primary string, secondary string, definitions []string, err error) {
	var builder strings.Builder
	var xbuilder strings.Builder
	prefix := "    "

	if slices.Contains(ignore, _type.Name()) {
		return "", "", ignore, nil
	}

	ignore = append(ignore, _type.Name())

	builder.WriteString(fmt.Sprintf("export type %s = {\n", _type.Name()))
	for i := _type.NumField() - 1; i >= 0; i-- {
		f := _type.Field(i)
		t := f.Type
		k := t.Kind()
		name := f.Name
		// We cannot use this, goja' runtime.ToValue() will ignore tags since it's not marshaling.
		//if tag := f.Tag.Get("json"); tag != "" {
		//	name = tag
		//}
		switch k {
		case
			reflect.Chan,
			reflect.Func,
			reflect.Pointer,
			reflect.Invalid,
			reflect.Interface,
			reflect.UnsafePointer:
			err = fmt.Errorf("type %s of kind %s is not supported", t.String(), k.String())
			return
		case
			reflect.Map:
			var primaryLoc string
			var secondaryLoc string
			if primaryLoc, secondaryLoc, definitions, err = Extract(t.Elem(), ignore); err != nil {
				return
			}
			builder.WriteString(fmt.Sprintf("%s%s: Record<string, %s>", prefix, name, t.Elem().Name()))
			if primaryLoc != "" {
				xbuilder.WriteString("\n")
				xbuilder.WriteString(primaryLoc)
			}
			if secondaryLoc != "" {
				xbuilder.WriteString("\n")
				xbuilder.WriteString(secondaryLoc)
			}
		case
			reflect.Slice,
			reflect.Array:
			var primaryLoc string
			var secondaryLoc string
			if primaryLoc, secondaryLoc, definitions, err = Extract(t.Elem(), ignore); err != nil {
				return
			}
			builder.WriteString(fmt.Sprintf("%s%s: %s[]", prefix, name, t.Elem().Name()))
			if primaryLoc != "" {
				xbuilder.WriteString("\n")
				xbuilder.WriteString(primaryLoc)
			}
			if secondaryLoc != "" {
				xbuilder.WriteString("\n")
				xbuilder.WriteString(secondaryLoc)
			}
		case
			reflect.Struct:
			var primaryLoc string
			var secondaryLoc string
			if primaryLoc, secondaryLoc, definitions, err = Extract(t, ignore); err != nil {
				return
			}
			builder.WriteString(fmt.Sprintf("%s%s: %s", prefix, name, t.Name()))
			if primaryLoc != "" {
				xbuilder.WriteString("\n")
				xbuilder.WriteString(primaryLoc)
			}
			if secondaryLoc != "" {
				xbuilder.WriteString("\n")
				xbuilder.WriteString(secondaryLoc)
			}
		case
			reflect.Bool:
			builder.WriteString(fmt.Sprintf("%s%s: boolean", prefix, name))
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
			builder.WriteString(fmt.Sprintf("%s%s: number", prefix, name))
		case
			reflect.String:
			builder.WriteString(fmt.Sprintf("%s%s: string", prefix, name))
		default:
			builder.WriteString(fmt.Sprintf("%s%s: unknown", prefix, name))
		}
		builder.WriteString("\n")
	}

	builder.WriteString("}\n")

	primary = builder.String()
	secondary = xbuilder.String()

	return
}
