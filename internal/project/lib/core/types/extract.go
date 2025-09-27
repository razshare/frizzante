package types

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func Extract(prefix string, _type reflect.Type, ignore []string) (primary string, secondary string, definitions []string, err error) {
	var builder strings.Builder
	var xbuilder strings.Builder

	definitions = ignore
	padding := "    "

	if slices.Contains(definitions, prefix+_type.Name()) {
		return
	}

	definitions = append(definitions, prefix+_type.Name())

	builder.WriteString(fmt.Sprintf("export type %s = {\n", prefix+_type.Name()))
	for i := _type.NumField() - 1; i >= 0; i-- {
		f := _type.Field(i)
		t := f.Type
		k := t.Kind()
		name := f.Name

		if strings.ToLower(f.Name[0:1]) == f.Name[0:1] {
			continue
		}

		if k == reflect.Pointer {
			t = t.Elem()
		}

		if t.Name() == "error" {
			builder.WriteString(fmt.Sprintf("%s%s: string\n", padding, name))
			continue
		}

		// We cannot use this, goja's runtime.ToValue() will ignore tags since it's not marshaling.
		//if tag := f.Tag.Get("json"); tag != "" {
		//	name = tag
		//}
		switch k {
		case
			reflect.Chan,
			reflect.Func,
			reflect.Invalid,
			reflect.UnsafePointer:
			err = fmt.Errorf("type %s of kind %s is not supported", t.String(), k.String())
			return
		case
			reflect.Map:
			var primaryLoc string
			var secondaryLoc string
			var definitionsLoc = make([]string, 0)
			if primaryLoc, secondaryLoc, definitionsLoc, err = Extract(prefix+name, t.Elem(), definitions); err != nil {
				return
			}
			definitions = append(definitions, definitionsLoc...)
			builder.WriteString(fmt.Sprintf("%s%s: Record<string, %s>", padding, name, t.Elem().Name()))
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
			var definitionsLoc = make([]string, 0)
			if primaryLoc, secondaryLoc, definitionsLoc, err = Extract(prefix+name, t.Elem(), definitions); err != nil {
				return
			}
			definitions = append(definitions, definitionsLoc...)
			builder.WriteString(fmt.Sprintf("%s%s: %s[]", padding, name, t.Elem().Name()))
			if primaryLoc != "" {
				xbuilder.WriteString("\n")
				xbuilder.WriteString(primaryLoc)
			}
			if secondaryLoc != "" {
				xbuilder.WriteString("\n")
				xbuilder.WriteString(secondaryLoc)
			}
		case
			reflect.Struct,
			reflect.Pointer,
			reflect.Interface:
			var primaryLoc string
			var secondaryLoc string
			var definitionsLoc = make([]string, 0)
			if primaryLoc, secondaryLoc, definitionsLoc, err = Extract(prefix+name, t, definitions); err != nil {
				return
			}
			definitions = append(definitions, definitionsLoc...)
			builder.WriteString(fmt.Sprintf("%s%s: %s", padding, name, t.Name()))
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
			builder.WriteString(fmt.Sprintf("%s%s: boolean", padding, name))
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
			builder.WriteString(fmt.Sprintf("%s%s: number", padding, name))
		case
			reflect.String:
			builder.WriteString(fmt.Sprintf("%s%s: string", padding, name))
		default:
			builder.WriteString(fmt.Sprintf("%s%s: unknown", padding, name))
		}
		builder.WriteString("\n")
	}

	builder.WriteString("}\n")

	primary = builder.String()
	secondary = xbuilder.String()

	return
}
