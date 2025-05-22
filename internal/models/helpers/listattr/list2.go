package listattr

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/listtype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type Type[T any] = listtype.ListNestedObjectValueOf[T]

func ValueOf[T any](ctx context.Context, values []*T) Type[T] {
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

func Default[T any](value []*T, attributes map[string]schema.Attribute, validators ...validator.Object) schema.ListNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.ListNestedAttribute{
		Optional:     true,
		Computed:     true,
		NestedObject: nested,
		CustomType:   listtype.NewListNestedObjectTypeOfMust[T](context.Background()),
		Default:      listdefault.StaticValue(ValueOf(context.Background(), value).ListValue),
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

	result := listtype.NewListNestedObjectValueOfSliceMust(h.Ctx, elems)

	// TODO
	h.Log("Setting list value for key '%s' of type '%T' to %s", key, result, types.UnsafeFormattedValue(result, true))
	*l = result
}
