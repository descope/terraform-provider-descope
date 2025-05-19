package types

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unsafe"
)

// Deprecated: UnsafeDebugValue is not safe to use in production code.
func UnsafeDebugValue(v any, options ...bool) string {
	var compact, types bool
	if len(options) > 0 {
		compact = options[0]
	}
	if len(options) > 1 {
		types = options[1]
	}
	if compact {
		return makeDebugValue(v, types, ", ", "", "")
	}
	return makeDebugValue(v, types, ",", " ", "\n")
}

func makeDebugValue(v any, types bool, sep, padding, endline string) string {
	var sb strings.Builder
	appendDebugValue(reflect.ValueOf(v), &sb, types, sep, padding, endline, 0)
	return sb.String()
}

func appendDebugValue(val reflect.Value, sb *strings.Builder, types bool, sep, padding, endline string, indent int) {
	typ := val.Type()
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() == reflect.Struct {
		sb.WriteString("{")
		for i := range val.NumField() {
			field := typ.Field(i)
			fieldVal := val.Field(i)
			if !fieldVal.CanInterface() {
				fieldVal = reflect.NewAt(fieldVal.Type(), unsafe.Pointer(fieldVal.UnsafeAddr())).Elem()
			}

			if i == 0 {
				sb.WriteString(endline)
				sb.WriteString(strings.Repeat(padding, indent+2))
			}
			if types {
				sb.WriteString(fmt.Sprintf("%s: %s = ", field.Name, field.Type))
			} else {
				sb.WriteString(fmt.Sprintf("%s: ", field.Name))
			}

			appendDebugValue(fieldVal, sb, types, sep, padding, endline, indent+2)

			if i < val.NumField()-1 {
				sb.WriteString(sep)
				sb.WriteString(endline)
				sb.WriteString(strings.Repeat(padding, indent+2))
			}
		}
		if val.NumField() > 0 {
			sb.WriteString(endline)
			sb.WriteString(strings.Repeat(padding, indent))
		}
		sb.WriteString("}")
	} else if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		sb.WriteString("[")
		for i := range val.Len() {
			if i == 0 {
				sb.WriteString(endline)
				sb.WriteString(strings.Repeat(padding, indent+2))
			}
			sb.WriteString(fmt.Sprintf("%d: ", i))

			appendDebugValue(val.Index(i), sb, types, sep, padding, endline, indent+2)

			if i < val.Len()-1 {
				sb.WriteString(sep)
				sb.WriteString(endline)
				sb.WriteString(strings.Repeat(padding, indent+2))
			}
		}
		if val.Len() > 0 {
			sb.WriteString(endline)
			sb.WriteString(strings.Repeat(padding, indent))
		}
		sb.WriteString("]")
	} else if val.Kind() == reflect.Map {
		m := map[string]string{}
		for _, k := range val.MapKeys() {
			s := strings.Builder{}
			appendDebugValue(k, &s, types, sep, padding, endline, indent+2)
			key := s.String()

			s.Reset()
			appendDebugValue(val.MapIndex(k), &s, types, sep, padding, endline, indent+2)
			value := s.String()

			m[key] = value
		}

		keys := []string{}
		for k := range m {
			keys = append(keys, k)
		}
		slices.Sort(keys)

		sb.WriteString("[")
		for i, key := range keys {
			if i == 0 {
				sb.WriteString(endline)
				sb.WriteString(strings.Repeat(padding, indent+2))
			}
			sb.WriteString(fmt.Sprintf("%s: %s", key, m[key]))
			if i < len(keys)-1 {
				sb.WriteString(sep)
				sb.WriteString(endline)
				sb.WriteString(strings.Repeat(padding, indent+2))
			}
		}
		if len(keys) > 0 {
			sb.WriteString(endline)
			sb.WriteString(strings.Repeat(padding, indent))
		}
		sb.WriteString("]")
	} else if val.Kind() == reflect.String {
		sb.WriteString(fmt.Sprintf("%q", val.Interface()))
	} else if val.CanInterface() {
		sb.WriteString(fmt.Sprint(val.Interface()))
	} else {
		sb.WriteString(fmt.Sprintf("<?%v?>", val))
	}
}
