package boolattr

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type = types.Bool

func Value(value bool) Type {
	return types.BoolValue(value)
}

func Required(validators ...validator.Bool) schema.BoolAttribute {
	return schema.BoolAttribute{
		Required:   true,
		Validators: validators,
	}
}

func Optional(validators ...validator.Bool) schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:      true,
		Computed:      true,
		Validators:    validators,
		PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
	}
}

func Default(value bool, validators ...validator.Bool) schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:   true,
		Computed:   true,
		Validators: validators,
		Default:    booldefault.StaticBool(value),
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
