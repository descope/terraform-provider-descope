package stringattr

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Identifier() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
}

func IdentifierMatched() schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,
	}
}

func Required(validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Required:   true,
		Validators: append([]validator.String{NonEmptyValidator}, validators...),
	}
}

func SecretRequired(validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Required:   true,
		Sensitive:  true,
		Validators: append([]validator.String{NonEmptyValidator}, validators...),
	}
}

func SecretOptional(validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:   true,
		Computed:   true,
		Sensitive:  true,
		Validators: validators,
		Default:    NullDefault(),
	}
}

func Optional(validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:      true,
		Computed:      true,
		Validators:    validators,
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
}

func Default(value string, validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:   true,
		Computed:   true,
		Validators: validators,
		Default:    stringdefault.StaticString(value),
	}
}

func Deprecated(message string, validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:           true,
		Computed:           true,
		DeprecationMessage: message + " This attribute will be removed in the next major version of the provider.",
		Validators:         validators,
		PlanModifiers:      []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		Default:            NullDefault(),
	}
}

func Renamed(oldname, newname string, validators ...validator.String) schema.StringAttribute {
	return Deprecated("The "+oldname+" attribute has been renamed, set the "+newname+" attribute instead.", validators...)
}

func Get(s types.String, data map[string]any, key string) {
	if !s.IsNull() && !s.IsUnknown() {
		data[key] = s.ValueString()
	}
}

func Set(s *types.String, data map[string]any, key string) {
	if v, ok := data[key].(string); ok {
		*s = types.StringValue(v)
	}
}

func EnsureKnown(s *types.String, defaultValue ...string) {
	if s.IsUnknown() {
		if len(defaultValue) > 0 {
			*s = types.StringValue(defaultValue[0])
		} else {
			*s = types.StringValue("")
		}
	}
}
