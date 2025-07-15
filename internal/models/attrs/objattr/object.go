package objattr

import (
	"context"
	"fmt"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/objtype"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Type[T any] = objtype.ObjectValueOf[T]

func Value[T any](value *T) Type[T] {
	return valueOf(context.Background(), value)
}

func valueOf[T any](ctx context.Context, value *T) Type[T] {
	if value == nil {
		return objtype.NewNullValue[T](ctx)
	}
	return helpers.Require(objtype.NewValue(ctx, value))
}

func Required[T any](attributes map[string]schema.Attribute, extras ...any) schema.SingleNestedAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.SingleNestedAttribute{
		Required:      true,
		CustomType:    objtype.NewType[T](context.Background()),
		Attributes:    attributes,
		Validators:    validators,
		PlanModifiers: modifiers,
	}
}

func Optional[T any](attributes map[string]schema.Attribute, extras ...any) schema.SingleNestedAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.SingleNestedAttribute{
		Optional:      true,
		Computed:      true,
		CustomType:    objtype.NewType[T](context.Background()),
		Attributes:    attributes,
		Validators:    validators,
		PlanModifiers: append([]planmodifier.Object{objectplanmodifier.UseStateForUnknown()}, modifiers...),
	}
}

func Default[T any](value *T, attributes map[string]schema.Attribute, extras ...any) schema.SingleNestedAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.SingleNestedAttribute{
		Optional:      true,
		Computed:      true,
		CustomType:    objtype.NewType[T](context.Background()),
		Attributes:    attributes,
		Validators:    validators,
		PlanModifiers: modifiers,
		Default:       objectdefault.StaticValue(Value(value).ObjectValue),
	}
}

func Get[T any, M helpers.Model[T]](o Type[T], data map[string]any, key string, h *helpers.Handler) {
	if o.IsUnknown() {
		return
	}
	if o.IsNull() {
		if key != helpers.RootKey {
			data[key] = nil
		}
		return
	}

	var value M = helpers.Require(o.ToObject(h.Ctx))
	if key == helpers.RootKey {
		maps.Copy(data, value.Values(h))
	} else if m, ok := data[key].(map[string]any); ok {
		maps.Copy(m, value.Values(h))
	} else {
		data[key] = value.Values(h)
	}
}

func Set[T any, M helpers.Model[T]](o *Type[T], data map[string]any, key string, h *helpers.Handler) {
	if !helpers.ShouldSetAttributeValue(h.Ctx, o) {
		return
	}

	var m map[string]any
	if key == helpers.RootKey {
		m = data
	} else if v, ok := data[key].(map[string]any); ok {
		m = v
	} else {
		*o = valueOf[T](h.Ctx, nil)
		return
	}

	var value M
	if o.IsNull() || o.IsUnknown() {
		value = new(T)
	} else {
		value = helpers.Require(o.ToObject(h.Ctx))
	}
	value.SetValues(h, m)

	*o = valueOf(h.Ctx, value)
}

func Nil[T any, M helpers.Model[T]](o *Type[T]) {
	if o.IsUnknown() {
		*o = Value[T](nil)
	}
}

func CollectReferences[T any, M helpers.CollectReferencesModel[T]](o Type[T], h *helpers.Handler) {
	if o.IsNull() || o.IsUnknown() {
		return
	}

	var value M = helpers.Require(o.ToObject(h.Ctx))
	value.CollectReferences(h)
}

func UpdateReferences[T any, M helpers.UpdateReferencesModel[T]](o *Type[T], h *helpers.Handler) {
	if o.IsNull() || o.IsUnknown() {
		return
	}

	var value M = helpers.Require(o.ToObject(h.Ctx))
	value.UpdateReferences(h)

	*o = valueOf(h.Ctx, value)
}

func parseExtras(extras []any) (validators []validator.Object, modifiers []planmodifier.Object) {
	for _, e := range extras {
		matched := false
		if validator, ok := e.(validator.Object); ok {
			matched = true
			validators = append(validators, validator)
		}
		if modifier, ok := e.(planmodifier.Object); ok {
			matched = true
			modifiers = append(modifiers, modifier)
		}
		if !matched {
			panic(fmt.Sprintf("unexpected extra value of type %T in object attribute", e))
		}
	}
	return
}

func modelFromObject[T any, M helpers.Model[T]](ctx context.Context, object types.Object, diagnostics *diag.Diagnostics) M {
	result := new(T)
	diags := object.As(ctx, result, basetypes.ObjectAsOptions{})
	diagnostics.Append(diags...)
	return result
}
