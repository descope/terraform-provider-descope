package intattr

import (
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type = types.Int64

func Value(value int64) Type {
	return types.Int64Value(value)
}

func Required(extras ...any) schema.Int64Attribute {
	validators, modifiers := parseExtras(extras)
	return schema.Int64Attribute{
		Required:      true,
		Validators:    validators,
		PlanModifiers: modifiers,
	}
}

func Optional(extras ...any) schema.Int64Attribute {
	validators, modifiers := parseExtras(extras)
	return schema.Int64Attribute{
		Optional:      true,
		Computed:      true,
		Validators:    validators,
		PlanModifiers: append([]planmodifier.Int64{helpers.UseValidStateForUnknown()}, modifiers...),
	}
}

func Default(value int, extras ...any) schema.Int64Attribute {
	validators, modifiers := parseExtras(extras)
	return schema.Int64Attribute{
		Optional:      true,
		Computed:      true,
		Validators:    validators,
		PlanModifiers: modifiers,
		Default:       int64default.StaticInt64(int64(value)),
	}
}

func Get(n types.Int64, data map[string]any, key string) {
	if !n.IsNull() && !n.IsUnknown() {
		data[key] = n.ValueInt64()
	}
}

func Set(n *types.Int64, data map[string]any, key string) {
	if v, ok := data[key].(float64); ok {
		*n = Value(int64(v))
	} else if v, ok := data[key].(int64); ok {
		*n = Value(v)
	} else if n.IsUnknown() {
		*n = Value(0)
	}
}

func parseExtras(extras []any) (validators []validator.Int64, modifiers []planmodifier.Int64) {
	for _, e := range extras {
		matched := false
		if validator, ok := e.(validator.Int64); ok {
			matched = true
			validators = append(validators, validator)
		}
		if modifier, ok := e.(planmodifier.Int64); ok {
			matched = true
			modifiers = append(modifiers, modifier)
		}
		if !matched {
			panic(fmt.Sprintf("unexpected extra value of type %T in int attribute", e))
		}
	}
	return
}
