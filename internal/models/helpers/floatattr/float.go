package floatattr

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Required(validators ...validator.Float64) schema.Float64Attribute {
	return schema.Float64Attribute{
		Required:   true,
		Validators: validators,
	}
}

func Optional(validators ...validator.Float64) schema.Float64Attribute {
	return schema.Float64Attribute{
		Optional:      true,
		Computed:      true,
		Validators:    validators,
		PlanModifiers: []planmodifier.Float64{float64planmodifier.UseStateForUnknown()},
	}
}

func Default(value float64, validators ...validator.Float64) schema.Float64Attribute {
	return schema.Float64Attribute{
		Optional:   true,
		Computed:   true,
		Validators: validators,
		Default:    float64default.StaticFloat64(value),
	}
}

func Get(n types.Float64, data map[string]any, key string) {
	if !n.IsNull() && !n.IsUnknown() {
		data[key] = n.ValueFloat64()
	}
}

func Set(n *types.Float64, data map[string]any, key string) {
	if v, ok := data[key].(float64); ok {
		*n = types.Float64Value(v)
	}
	if v, ok := data[key].(int64); ok {
		*n = types.Float64Value(float64(v))
	}
}
