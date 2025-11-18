package mapattr

import (
	"context"
	"fmt"
	"iter"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/maptype"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/objtype"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type Type[T any] = maptype.MapValueOf[T]

func Value[T any](value map[string]*T) Type[T] {
	return valueOf(context.Background(), value)
}

func Empty[T any]() Type[T] {
	return valueOf(context.Background(), map[string]*T{})
}

func valueOf[T any](ctx context.Context, value map[string]*T) Type[T] {
	if value == nil {
		return maptype.NewNullValue[T](ctx)
	}
	return helpers.Require(maptype.NewValue(ctx, value))
}

func Required[T any](attributes map[string]schema.Attribute, extras ...any) schema.MapNestedAttribute {
	mapValidators, objectValidators := parseExtras(extras)
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: objectValidators,
	}
	return schema.MapNestedAttribute{
		Required:     true,
		NestedObject: nested,
		CustomType:   maptype.NewType[T](context.Background()),
		Validators:   mapValidators,
	}
}

func Optional[T any](attributes map[string]schema.Attribute, extras ...any) schema.MapNestedAttribute {
	mapValidators, objectValidators := parseExtras(extras)
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: objectValidators,
	}
	return schema.MapNestedAttribute{
		Optional:      true,
		Computed:      true,
		NestedObject:  nested,
		CustomType:    maptype.NewType[T](context.Background()),
		PlanModifiers: []planmodifier.Map{helpers.UseValidStateForUnknown()},
		Validators:    mapValidators,
	}
}

func Default[T any](values map[string]*T, attributes map[string]schema.Attribute, extras ...any) schema.MapNestedAttribute {
	mapValidators, objectValidators := parseExtras(extras)
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: objectValidators,
	}
	return schema.MapNestedAttribute{
		Optional:     true,
		Computed:     true,
		NestedObject: nested,
		CustomType:   maptype.NewType[T](context.Background()),
		Default:      mapdefault.StaticValue(Value(values).MapValue),
		Validators:   mapValidators,
	}
}

func Get[T any, M helpers.Model[T]](m Type[T], data map[string]any, key string, h *helpers.Handler) {
	if m.IsUnknown() {
		return
	}
	if m.IsNull() {
		if key != helpers.RootKey {
			data[key] = nil
		}
		return
	}

	elems, diags := m.ToMap(h.Ctx)
	h.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	result := map[string]any{}
	for k, v := range elems {
		var element M = v
		result[k] = element.Values(h)
	}

	data[key] = result
}

func Set[T any, M helpers.Model[T]](m *Type[T], data map[string]any, key string, h *helpers.Handler) {
	if !helpers.ShouldSetAttributeValue(h.Ctx, m) {
		return
	}

	values := data
	if key != helpers.RootKey {
		values, _ = data[key].(map[string]any)
	}

	elems := map[string]*T{}
	current := m.Elements()

	for k, v := range values {
		var element M
		if c, ok := current[k]; ok && !c.IsNull() && !c.IsUnknown() {
			element, _ = objtype.NewObjectWith[T](h.Ctx, c)
		}
		if element == nil {
			element = new(T)
		}
		if modelData, ok := v.(map[string]any); ok {
			element.SetValues(h, modelData)
		}
		elems[k] = element
	}

	*m = valueOf(h.Ctx, elems)
}

func Iterator[T any](m Type[T], h *helpers.Handler) iter.Seq2[string, *T] {
	return func(yield func(string, *T) bool) {
		for k, v := range m.Elements() {
			if v.IsNull() || v.IsUnknown() {
				continue
			}

			ptr, diags := objtype.NewObjectWith[T](h.Ctx, v)
			h.Diagnostics.Append(diags...)
			if diags.HasError() {
				continue
			}

			if !yield(k, ptr) {
				break
			}
		}
	}
}

func MutatingIterator[T any](m *Type[T], h *helpers.Handler) iter.Seq2[string, *T] {
	return func(yield func(string, *T) bool) {
		elements := m.Elements()

		for k, v := range elements {
			if v.IsNull() || v.IsUnknown() {
				continue
			}

			ptr, diags := objtype.NewObjectWith[T](h.Ctx, v)
			h.Diagnostics.Append(diags...)
			if diags.HasError() {
				continue
			}

			cont := yield(k, ptr)

			obj, diags := objtype.NewValue(h.Ctx, ptr)
			h.Diagnostics.Append(diags...)
			if !diags.HasError() {
				elements[k] = obj
			}

			if !cont {
				break
			}
		}

		mapValue, diags := maptype.NewValueWith[T](h.Ctx, elements)
		h.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		*m = mapValue
	}
}

func parseExtras(extras []any) (mapValidators []validator.Map, objectValidators []validator.Object) {
	for _, e := range extras {
		matched := false
		if v, ok := e.(validator.Map); ok {
			matched = true
			mapValidators = append(mapValidators, v)
		}
		if v, ok := e.(validator.Object); ok {
			matched = true
			objectValidators = append(objectValidators, v)
		}
		if !matched {
			panic(fmt.Sprintf("unexpected extra value of type %T in attribute", e))
		}
	}
	return
}
