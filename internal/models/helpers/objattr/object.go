package objattr

import (
	"context"
	"fmt"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objattr/objtype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type WithType[T any] = objtype.ObjectValueOf[T]

func Required[T any](attributes map[string]schema.Attribute, extras ...any) schema.SingleNestedAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.SingleNestedAttribute{
		Required:      true,
		CustomType:    objtype.NewObjectTypeOf[T](context.Background()),
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
		CustomType:    objtype.NewObjectTypeOf[T](context.Background()),
		Attributes:    attributes,
		Validators:    validators,
		PlanModifiers: append([]planmodifier.Object{objectplanmodifier.UseStateForUnknown()}, modifiers...),
	}
}

func Get[T any, M helpers.Model[T]](o objtype.ObjectValueOf[T], data map[string]any, key string, h *helpers.Handler) {
	if o.IsNull() || o.IsUnknown() {
		return
	}

	var v M = o.ToPtrMust(h.Ctx)
	if m, ok := data[key].(map[string]any); ok {
		maps.Copy(m, v.Values(h))
	} else {
		data[key] = v.Values(h)
	}
}

func Set[T any, M helpers.Model[T]](o *objtype.ObjectValueOf[T], data map[string]any, key string, h *helpers.Handler) {
	m, ok := data[key].(map[string]any)
	if !ok {
		return
	}

	var v M = new(T)
	v.SetValues(h, m)

	*o = objtype.NewObjectValueOfMust(h.Ctx, v)
}

func Ensure[T any, M helpers.Model[T]](o *objtype.ObjectValueOf[T], data map[string]any, key string, h *helpers.Handler) {
	if o.IsUnknown() {
		Set[T, M](o, data, key, h)
	}
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
