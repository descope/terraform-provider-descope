package objattr

import (
	"context"
	"fmt"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type Type[T any] = objtype.ObjectValueOf[T]

func ValueOf[T any](ctx context.Context, value *T) Type[T] {
	return objtype.NewObjectValueOfMust(ctx, value)
}

func Required[T any](attributes map[string]schema.Attribute, extras ...any) schema.SingleNestedAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.SingleNestedAttribute{
		Required:      true,
		CustomType:    objtype.NewObjectTypeOfMust[T](context.Background()),
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
		CustomType:    objtype.NewObjectTypeOfMust[T](context.Background()),
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
		CustomType:    objtype.NewObjectTypeOfMust[T](context.Background()),
		Attributes:    attributes,
		Validators:    validators,
		PlanModifiers: modifiers,
		Default:       objectdefault.StaticValue(ValueOf(context.Background(), value).ObjectValue),
	}
}

func Get[T any, M helpers.Model[T]](o Type[T], data map[string]any, key string, h *helpers.Handler) {
	if o.IsNull() || o.IsUnknown() {
		return
	}

	var value M = o.ToPtrMust(h.Ctx)
	if key == helpers.RootKey {
		maps.Copy(data, value.Values(h))
	} else if m, ok := data[key].(map[string]any); ok {
		maps.Copy(m, value.Values(h))
	} else {
		data[key] = value.Values(h)
	}
}

func Set[T any, M helpers.Model[T]](o *Type[T], data map[string]any, key string, h *helpers.Handler) {
	var m map[string]any
	if key == helpers.RootKey {
		m = data
	} else if v, ok := data[key].(map[string]any); ok {
		m = v
	} else {
		return
	}

	var value M
	if o.IsNull() || o.IsUnknown() {
		value = new(T)
	} else {
		value = o.ToPtrMust(h.Ctx)
	}
	value.SetValues(h, m)

	// TODO
	result := objtype.NewObjectValueOfMust(h.Ctx, value)
	h.Log("Setting object value for key '%s' of type '%T' to %s", key, result, types.UnsafeFormattedValue(result, true))
	*o = result
}

func CollectReferences[T any, M helpers.CollectReferencesModel[T]](o Type[T], h *helpers.Handler) {
	if o.IsNull() || o.IsUnknown() {
		return
	}

	var value M = o.ToPtrMust(h.Ctx)
	value.CollectReferences(h)
}

func UpdateReferences[T any, M helpers.UpdateReferencesModel[T]](o *Type[T], h *helpers.Handler) {
	if o.IsNull() || o.IsUnknown() {
		return
	}

	var value M = o.ToPtrMust(h.Ctx)
	value.UpdateReferences(h)

	*o = objtype.NewObjectValueOfMust(h.Ctx, value)
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
