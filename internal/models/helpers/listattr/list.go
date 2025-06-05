package listattr

import (
	"context"
	"iter"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/listtype"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type Type[T any] = listtype.ListValueOf[T]

func Value[T any](values []*T) Type[T] {
	return valueOf(context.Background(), values)
}

func Empty[T any]() Type[T] {
	return valueOf(context.Background(), []*T{})
}

func valueOf[T any](ctx context.Context, values []*T) Type[T] {
	return helpers.Require(listtype.NewValue(ctx, values))
}

func Required[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.ListNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.ListNestedAttribute{
		Required:     true,
		NestedObject: nested,
		CustomType:   listtype.NewType[T](context.Background()),
	}
}

func Optional[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.ListNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.ListNestedAttribute{
		Optional:      true,
		Computed:      true,
		NestedObject:  nested,
		CustomType:    listtype.NewType[T](context.Background()),
		PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
	}
}

func Default[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.ListNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.ListNestedAttribute{
		Optional:     true,
		Computed:     true,
		NestedObject: nested,
		CustomType:   listtype.NewType[T](context.Background()),
		Default:      listdefault.StaticValue(Empty[T]().ListValue),
	}
}

func Get[T any, M helpers.Model[T]](l Type[T], data map[string]any, key string, h *helpers.Handler) {
	if l.IsNull() || l.IsUnknown() {
		return
	}

	elems, diags := l.ToSlice(h.Ctx)
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

func Set[T any, M helpers.Model[T]](l *Type[T], data map[string]any, key string, h *helpers.Handler) {
	values, _ := data[key].([]any)

	elems := []*T{}
	current := l.Elements()

	for i, v := range values {
		var element M
		if len(current) > i && !current[i].IsNull() && !current[i].IsUnknown() {
			element, _ = objtype.NewObjectWith[T](h.Ctx, current[i])
		}
		if element == nil {
			element = new(T)
		}
		if modelData, ok := v.(map[string]any); ok {
			element.SetValues(h, modelData)
		}
		elems = append(elems, element)
	}

	*l = valueOf(h.Ctx, elems)
}

func Iterator[T any](l Type[T], h *helpers.Handler) iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for _, v := range l.Elements() {
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

func MutatingIterator[T any](l *Type[T], h *helpers.Handler) iter.Seq[*T] {
	return func(yield func(*T) bool) {
		elements := l.Elements()

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

		listValue, diags := listtype.NewValueWith[T](h.Ctx, elements)
		h.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		*l = listValue
	}
}
