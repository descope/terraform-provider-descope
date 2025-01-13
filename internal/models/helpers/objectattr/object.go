package objectattr

import (
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Required(attributes map[string]schema.Attribute, extras ...any) schema.SingleNestedAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.SingleNestedAttribute{
		Required:      true,
		Attributes:    attributes,
		Validators:    validators,
		PlanModifiers: modifiers,
	}
}

func Optional(attributes map[string]schema.Attribute, extras ...any) schema.SingleNestedAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.SingleNestedAttribute{
		Optional:      true,
		Computed:      true,
		Attributes:    attributes,
		Validators:    validators,
		PlanModifiers: modifiers,
		Default:       objectdefault.StaticValue(types.ObjectNull(getAttributeTypes(attributes))),
	}
}

func Get[T any, M helpers.Model[T]](o M, data map[string]any, key string, h *helpers.Handler) {
	if o != nil {
		if o, ok := any(o).(checkableModel); ok {
			o.Check(h)
		}
		data[key] = o.Values(h)
	}
}

func Set[T any, M helpers.Model[T]](o *M, data map[string]any, key string, h *helpers.Handler) { // using *M to stay consistent with other Set functions
	if v, ok := data[key].(map[string]any); ok {
		if *o != nil {
			(*o).SetValues(h, v)
		}
	}
}

type checkableModel interface {
	Check(*helpers.Handler)
}

func getAttributeTypes(attributes map[string]schema.Attribute) map[string]attr.Type {
	types := map[string]attr.Type{}
	for k, v := range attributes {
		types[k] = v.GetType()
	}
	return types
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
