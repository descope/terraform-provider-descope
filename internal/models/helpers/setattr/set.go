package setattr

import (
	"context"
	"iter"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/settype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type Type[T any] = settype.SetNestedObjectValueOf[T]

func ValueOf[T any](ctx context.Context, values []*T) Type[T] {
	return settype.NewSetNestedObjectValueOfSliceMust(ctx, values)
}

func Required[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.SetNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.SetNestedAttribute{
		Required:     true,
		NestedObject: nested,
		CustomType:   settype.NewSetNestedObjectTypeOfMust[T](context.Background()),
	}
}

func Optional[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.SetNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.SetNestedAttribute{
		Optional:      true,
		Computed:      true,
		NestedObject:  nested,
		CustomType:    settype.NewSetNestedObjectTypeOfMust[T](context.Background()),
		PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
	}
}

func Default[T any](value []*T, attributes map[string]schema.Attribute, validators ...validator.Object) schema.SetNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.SetNestedAttribute{
		Optional:     true,
		Computed:     true,
		NestedObject: nested,
		CustomType:   settype.NewSetNestedObjectTypeOfMust[T](context.Background()),
		Default:      setdefault.StaticValue(ValueOf(context.Background(), value).SetValue),
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

	result := settype.NewSetNestedObjectValueOfSliceMust(h.Ctx, elems)

	// TODO
	h.Log("Setting set value for key '%s' of type '%T' to %s", key, result, types.UnsafeFormattedValue(result, true))
	*l = result
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
