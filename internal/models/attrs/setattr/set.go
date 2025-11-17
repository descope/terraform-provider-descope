package setattr

import (
	"context"
	"iter"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/objtype"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/settype"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type Type[T any] = settype.SetValueOf[T]

func Value[T any](values []*T) Type[T] {
	return valueOf(context.Background(), values)
}

func Empty[T any]() Type[T] {
	return valueOf(context.Background(), []*T{})
}

func valueOf[T any](ctx context.Context, values []*T) Type[T] {
	return helpers.Require(settype.NewValue(ctx, values))
}

// Deprecated: The set type is buggy, use a list instead.
func Required[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.SetNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.SetNestedAttribute{
		Required:     true,
		NestedObject: nested,
		CustomType:   settype.NewType[T](context.Background()),
	}
}

// Deprecated: The set type is buggy, use a list instead.
func Optional[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.SetNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.SetNestedAttribute{
		Optional:      true,
		Computed:      true,
		NestedObject:  nested,
		CustomType:    settype.NewType[T](context.Background()),
		PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
	}
}

// Deprecated: The set type is buggy, use a list instead.
func Default[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.SetNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.SetNestedAttribute{
		Optional:     true,
		Computed:     true,
		NestedObject: nested,
		CustomType:   settype.NewType[T](context.Background()),
		Default:      setdefault.StaticValue(Empty[T]().SetValue),
	}
}

func Get[T any, M helpers.Model[T]](s Type[T], data map[string]any, key string, h *helpers.Handler) {
	if s.IsNull() || s.IsUnknown() {
		return
	}

	elems, diags := s.ToSlice(h.Ctx)
	h.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	result := []any{}
	for _, v := range elems {
		var m M = v
		result = append(result, m.Values(h))
	}

	data[key] = result
}

func Set[T any, M helpers.Model[T]](s *Type[T], data map[string]any, key string, h *helpers.Handler) {
	values, _ := data[key].([]any)

	elems := []*T{}

	for _, v := range values {
		var element M = new(T)
		if modelData, ok := v.(map[string]any); ok {
			element.SetValues(h, modelData)
		}
		elems = append(elems, element)
	}

	*s = valueOf(h.Ctx, elems)
}

func Iterator[T any](s Type[T], h *helpers.Handler) iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for _, v := range s.Elements() {
			if v.IsNull() || v.IsUnknown() {
				continue
			}

			ptr, diags := objtype.NewObjectWith[T](h.Ctx, v)
			h.Diagnostics.Append(diags...)
			if diags.HasError() {
				continue
			}

			if !yield(ptr) {
				break
			}
		}
	}
}

func MutatingIterator[T any](s *Type[T], h *helpers.Handler) iter.Seq[*T] {
	return func(yield func(*T) bool) {
		elements := s.Elements()

		for i, v := range elements {
			if v.IsNull() || v.IsUnknown() {
				continue
			}

			ptr, diags := objtype.NewObjectWith[T](h.Ctx, v)
			h.Diagnostics.Append(diags...)
			if diags.HasError() {
				continue
			}

			cont := yield(ptr)

			obj, diags := objtype.NewValue(h.Ctx, ptr)
			h.Diagnostics.Append(diags...)
			if !diags.HasError() {
				elements[i] = obj
			}

			if !cont {
				break
			}
		}

		setValue, diags := settype.NewValueWith[T](h.Ctx, elements)
		h.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		*s = setValue
	}
}
