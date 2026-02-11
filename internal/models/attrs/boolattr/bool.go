package boolattr

import (
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type = types.Bool

func Value(value bool) Type {
	return types.BoolValue(value)
}

func Required(extras ...any) schema.BoolAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.BoolAttribute{
		Required:      true,
		Validators:    validators,
		PlanModifiers: modifiers,
	}
}

func Optional(extras ...any) schema.BoolAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.BoolAttribute{
		Optional:      true,
		Computed:      true,
		Validators:    validators,
		PlanModifiers: append([]planmodifier.Bool{helpers.UseValidStateForUnknown()}, modifiers...),
	}
}

func Default(value bool, extras ...any) schema.BoolAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.BoolAttribute{
		Optional:      true,
		Computed:      true,
		Validators:    validators,
		PlanModifiers: modifiers,
		Default:       booldefault.StaticBool(value),
	}
}

func Get(b types.Bool, data map[string]any, key string) {
	if !b.IsNull() && !b.IsUnknown() {
		data[key] = b.ValueBool()
	}
}

func Set(b *types.Bool, data map[string]any, key string) {
	if v, ok := data[key].(bool); ok {
		*b = Value(v)
	} else if b.IsUnknown() {
		*b = Value(false)
	}
}

func GetNot(b types.Bool, data map[string]any, key string) {
	if !b.IsNull() && !b.IsUnknown() {
		data[key] = !b.ValueBool()
	}
}

func SetNot(b *types.Bool, data map[string]any, key string) {
	if v, ok := data[key].(bool); ok {
		*b = Value(!v)
	} else if b.IsUnknown() {
		*b = Value(true)
	}
}

func parseExtras(extras []any) (validators []validator.Bool, modifiers []planmodifier.Bool) {
	for _, e := range extras {
		matched := false
		if validator, ok := e.(validator.Bool); ok {
			matched = true
			validators = append(validators, validator)
		}
		if modifier, ok := e.(planmodifier.Bool); ok {
			matched = true
			modifiers = append(modifiers, modifier)
		}
		if !matched {
			panic(fmt.Sprintf("unexpected extra value of type %T in bool attribute", e))
		}
	}
	return
}
