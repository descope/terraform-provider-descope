package helpers

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unsafe"
)

// Deprecated: UnsafeFormattedValue is not safe to use in production code.
func UnsafeFormattedValue[T any](v T, pretty bool) string {
	f := UnsafeFormatter[T]{Private: true, Separator: ", "}
	if pretty {
		f = UnsafeFormatter[T]{Private: true, Separator: ",", Padding: " ", Endline: "\n", Indent: 4}
	}
	return f.Format(v)
}

// Deprecated: UnsafeFormatter is not safe to use in production code.
type UnsafeFormatter[T any] struct {
	Types     bool
	Private   bool
	Separator string
	Padding   string
	Endline   string
	Indent    int
}

func (f UnsafeFormatter[T]) Format(v T) string {
	var w strings.Builder

	val := reflect.ValueOf(v)
	if f.Types {
		w.WriteString(val.Type().Name())
	}
	if !val.CanAddr() {
		val = reflect.ValueOf(&v)
	}

	f.append(val, 0, &w)
	return w.String()
}

func (f UnsafeFormatter[T]) append(val reflect.Value, nesting int, b *strings.Builder) {
	typ := val.Type()
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if f.Private && !val.CanInterface() && val.CanAddr() {
		val = reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()
	}

	if val.Kind() == reflect.Struct {
		fields := []int{}
		for i := range val.NumField() {
			if f.Private || val.Field(i).CanInterface() {
				fields = append(fields, i)
			}
		}

		f.open("{", len(fields), nesting, b)
		for _, idx := range fields {
			field := typ.Field(idx)
			if f.Types {
				b.WriteString(fmt.Sprintf("%s: %s = ", field.Name, field.Type))
			} else {
				b.WriteString(fmt.Sprintf("%s: ", field.Name))
			}
			f.append(val.Field(idx), nesting+f.Indent, b)
			f.separate(idx, len(fields), nesting, b)
		}
		f.close("}", len(fields), nesting, b)
	} else if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		f.open("[", val.Len(), nesting, b)
		for i := range val.Len() {
			b.WriteString(fmt.Sprintf("%d: ", i))
			f.append(val.Index(i), nesting+f.Indent, b)
			f.separate(i, val.Len(), nesting, b)
		}
		f.close("]", val.Len(), nesting, b)
	} else if val.Kind() == reflect.Map {
		keys, values := f.convertMap(val, nesting)
		f.open("[", len(keys), nesting, b)
		for i := range keys {
			b.WriteString(fmt.Sprintf("%s: %s", keys[i], values[i]))
			f.separate(i, len(keys), nesting, b)
		}
		f.close("]", len(keys), nesting, b)
	} else if !val.CanInterface() {
		b.WriteString(fmt.Sprintf("<!%v!>", val))
	} else if val.Kind() == reflect.String {
		b.WriteString(fmt.Sprintf("%q", val.Interface()))
	} else {
		if f.Types && nesting == 0 {
			b.WriteString("(")
		}
		b.WriteString(fmt.Sprint(val.Interface()))
		if f.Types && nesting == 0 {
			b.WriteString(")")
		}
	}
}

func (f UnsafeFormatter[T]) open(s string, n int, nesting int, b *strings.Builder) {
	b.WriteString(s)
	if n > 0 {
		b.WriteString(f.Endline)
		b.WriteString(strings.Repeat(f.Padding, nesting+f.Indent))
	}
}

func (f UnsafeFormatter[T]) separate(i int, n int, nesting int, b *strings.Builder) {
	if i < n-1 {
		b.WriteString(f.Separator)
		b.WriteString(f.Endline)
		b.WriteString(strings.Repeat(f.Padding, nesting+f.Indent))
	}
}

func (f UnsafeFormatter[T]) close(s string, n int, nesting int, b *strings.Builder) {
	if n > 0 {
		b.WriteString(f.Endline)
		b.WriteString(strings.Repeat(f.Padding, nesting))
	}
	b.WriteString(s)
}

func (f UnsafeFormatter[T]) convertMap(val reflect.Value, nesting int) (keys []string, values []string) {
	m := map[string]string{}
	for _, k := range val.MapKeys() {
		key := strings.Builder{}
		f.append(k, nesting+f.Indent, &key)

		value := strings.Builder{}
		f.append(val.MapIndex(k), nesting+f.Indent, &value)

		m[key.String()] = value.String()
	}

	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	for _, key := range keys {
		values = append(values, m[key])
	}

	return keys, values
}
