package mapattr

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/maptype"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type Type[T any] = maptype.MapNestedObjectValueOf[T]

func ValueOf[T any](ctx context.Context, value map[string]*T) Type[T] {
	return maptype.NewMapNestedObjectValueOfMapMust(ctx, value)
}

func Required2[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.MapNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.MapNestedAttribute{
		Required:     true,
		NestedObject: nested,
		CustomType:   maptype.NewMapNestedObjectTypeOfMust[T](context.Background()),
	}
}

func Optional2[T any](attributes map[string]schema.Attribute, validators ...validator.Object) schema.MapNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.MapNestedAttribute{
		Optional:      true,
		Computed:      true,
		NestedObject:  nested,
		CustomType:    maptype.NewMapNestedObjectTypeOfMust[T](context.Background()),
		PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
	}
}

func Default[T any](values map[string]*T, attributes map[string]schema.Attribute, validators ...validator.Object) schema.MapNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.MapNestedAttribute{
		Optional:     true,
		Computed:     true,
		NestedObject: nested,
		CustomType:   maptype.NewMapNestedObjectTypeOfMust[T](context.Background()),
		Default:      mapdefault.StaticValue(ValueOf(context.Background(), values).MapValue),
	}
}

func Get2[T any, M helpers.Model[T]](m Type[T], data map[string]any, key string, h *helpers.Handler) {
	if m.IsNull() || m.IsUnknown() {
		return
	}

	elems, diags := m.ToMap(h.Ctx)
	h.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	result := map[string]any{}
	for k, v := range elems {
		var m M = v
		result[k] = m.Values(h)
	}

	data[key] = result
}

func Set2[T any, M helpers.Model[T]](m *Type[T], data map[string]any, key string, h *helpers.Handler) {
	values := data
	if key != helpers.RootKey {
		values, _ = data[key].(map[string]any)
	}

	elems := map[string]*T{}
	current := m.Elements()

	for k, v := range values {
		if modelData, ok := v.(map[string]any); ok {
			var element M
			if c := current[k]; c.IsNull() || c.IsUnknown() {
				element = new(T)
			} else {
				element = objtype.ObjectValueObjectPtrMust[T](h.Ctx, c)
			}
			element.SetValues(h, modelData)
			elems[k] = element
		}
	}

	result := maptype.NewMapNestedObjectValueOfMapMust(h.Ctx, elems)

	// TODO
	h.Log("Setting map value for key '%s' of type '%T' to %s", key, result, types.UnsafeFormattedValue(result, true))
	*m = result
}
