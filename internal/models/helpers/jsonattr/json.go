package jsonattr

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Required(validators ...validator.Dynamic) schema.DynamicAttribute {
	return schema.DynamicAttribute{
		Required:   true,
		Validators: append([]validator.Dynamic{typeValidator}, validators...),
	}
}

func Optional(validators ...validator.Dynamic) schema.DynamicAttribute {
	return schema.DynamicAttribute{
		Optional:   true,
		Computed:   true,
		Validators: append([]validator.Dynamic{typeValidator}, validators...),
		Default:    dynamicdefault.StaticValue(types.DynamicNull()),
	}
}

func Get(s types.Dynamic, data map[string]any, key string) {
	if !s.IsNull() && !s.IsUnknown() && !s.IsUnderlyingValueNull() && !s.IsUnderlyingValueUnknown() {
		if object, ok := s.UnderlyingValue().(types.Object); ok {
			data[key] = convertValue(object)
		}
	}
}

func convertValue(value attr.Value) any {
	switch v := value.(type) {
	case types.String:
		return v.ValueString()
	case types.Bool:
		return v.ValueBool()
	case types.Number:
		n, _ := v.ValueBigFloat().Float64()
		return n
	case types.Object:
		m := map[string]any{}
		for k, e := range v.Attributes() {
			m[k] = convertValue(e)
		}
		return m
	case types.Tuple:
		s := make([]any, len(v.Elements()))
		for i, e := range v.Elements() {
			s[i] = convertValue(e)
		}
		return s
	default:
		panic("unexpected type in JSON object")
	}
}
