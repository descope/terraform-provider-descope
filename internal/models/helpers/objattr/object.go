package objattr

import (
	"context"
	"fmt"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type Type[T any] = objtype.ObjectValueOf[T]

func ValueOf[T any](value *T) Type[T] {
	return objtype.NewObjectValueOfMust(context.Background(), value)
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

func Default[T any](attributes map[string]schema.Attribute, value *T, extras ...any) schema.SingleNestedAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.SingleNestedAttribute{
		Optional:      true,
		Computed:      true,
		CustomType:    objtype.NewObjectTypeOfMust[T](context.Background()),
		Attributes:    attributes,
		Validators:    validators,
		PlanModifiers: modifiers,
		Default:       objectdefault.StaticValue(ValueOf(value).ObjectValue),
	}
}

func Get[T any, M helpers.Model[T]](o Type[T], data map[string]any, key string, h *helpers.Handler) {
	if o.IsNull() || o.IsUnknown() {
		return
	}

	var value M = o.ToPtrMust(h.Ctx)
	if m, ok := data[key].(map[string]any); ok {
		maps.Copy(m, value.Values(h))
	} else {
		data[key] = value.Values(h)
	}
}

func Set[T any, M helpers.Model[T]](o *Type[T], data map[string]any, key string, h *helpers.Handler) {
	m, ok := data[key].(map[string]any)
	if !ok {
		return
	}

	var value M
	if o.IsNull() || o.IsUnknown() {
		value = new(T)
	} else {
		value = o.ToPtrMust(h.Ctx)
	}
	value.SetValues(h, m)

	*o = ValueOf(value)
}

func Ensure[T any, M helpers.Model[T]](o *Type[T], data map[string]any, key string, h *helpers.Handler) {
	if o.IsUnknown() {
		Set[T, M](o, data, key, h)
	}
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

	*o = ValueOf(value)
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
