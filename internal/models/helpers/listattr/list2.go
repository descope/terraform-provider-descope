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

type Type[T any] = listtype.ListNestedObjectValueOf[T]

func Value[T any](values []*T) Type[T] {
	return valueOf(context.Background(), values)
}

func Empty[T any]() Type[T] {
	return valueOf(context.Background(), []*T{})
}

func valueOf[T any](ctx context.Context, values []*T) Type[T] {
	return listtype.NewListNestedObjectValueOfSliceMust(ctx, values)
}

func Required2[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.ListNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.ListNestedAttribute{
		Required:     true,
		NestedObject: nested,
		CustomType:   listtype.NewListNestedObjectTypeOfMust[T](context.Background()),
	}
}

func Optional2[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.ListNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.ListNestedAttribute{
		Optional:      true,
		Computed:      true,
		NestedObject:  nested,
		CustomType:    listtype.NewListNestedObjectTypeOfMust[T](context.Background()),
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
		CustomType:   listtype.NewListNestedObjectTypeOfMust[T](context.Background()),
		Default:      listdefault.StaticValue(Empty[T]().ListValue),
	}
}

func Get2[T any, M helpers.Model[T]](l Type[T], data map[string]any, key string, h *helpers.Handler) {
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

func Set2[T any, M helpers.Model[T]](l *Type[T], data map[string]any, key string, h *helpers.Handler) {
	elems := []*T{}

	values, _ := data[key].([]any)
	for _, v := range values {
		if modelData, ok := v.(map[string]any); ok {
			var m M
			model := &m
			*model = new(T)
			(*model).SetValues(h, modelData)
			elems = append(elems, *model)
		}
	}

	*l = valueOf(h.Ctx, elems)
}

func Iterator[T any, M helpers.Model[T]](l Type[T], h *helpers.Handler) iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for _, v := range l.Elements() {
			if v.IsNull() || v.IsUnknown() {
				continue
			}

			ptr, diags := objtype.ObjectValueObjectPtr[T](h.Ctx, v)
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

func MutatingIterator[T any, M helpers.Model[T]](l *Type[T], h *helpers.Handler) iter.Seq[*T] {
	return func(yield func(*T) bool) {
		elements := l.Elements()

		for i, v := range elements {
			if v.IsNull() || v.IsUnknown() {
				continue
			}

			ptr, diags := objtype.ObjectValueObjectPtr[T](h.Ctx, v)
			h.Diagnostics.Append(diags...)
			if diags.HasError() {
				continue
			}

			cont := yield(ptr)

			obj, diags := objtype.NewObjectValueOf(h.Ctx, ptr)
			h.Diagnostics.Append(diags...)
			if !diags.HasError() {
				elements[i] = obj
			}

			if !cont {
				break
			}
		}

		listValue, diags := listtype.ValueOf[T](h.Ctx, elements)
		h.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		*l = listValue
	}
}
