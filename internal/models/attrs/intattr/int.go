package intattr

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type = types.Int64

func Value(value int64) Type {
	return types.Int64Value(value)
}

func Required(validators ...validator.Int64) schema.Int64Attribute {
	return schema.Int64Attribute{
		Required:   true,
		Validators: validators,
	}
}

func Optional(validators ...validator.Int64) schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:      true,
		Computed:      true,
		Validators:    validators,
		PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
	}
}

func Default(value int, validators ...validator.Int64) schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:   true,
		Computed:   true,
		Validators: validators,
		Default:    int64default.StaticInt64(int64(value)),
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
