package stringattr

import (
	"fmt"
	"slices"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type = types.String

func Value(value string) Type {
	return types.StringValue(value)
}

func Identifier() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:      true,
		PlanModifiers: []planmodifier.String{helpers.UseValidStateForUnknown()},
	}
}

func IdentifierMatched() schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,
	}
}

func Required(extras ...any) schema.StringAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.StringAttribute{
		Required:      true,
		Validators:    append([]validator.String{NonEmptyValidator}, validators...),
		PlanModifiers: modifiers,
	}
}

func SecretRequired(extras ...any) schema.StringAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.StringAttribute{
		Required:      true,
		Sensitive:     true,
		Validators:    append([]validator.String{NonEmptyValidator}, validators...),
		PlanModifiers: modifiers,
	}
}

func SecretOptional(extras ...any) schema.StringAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.StringAttribute{
		Optional:      true,
		Computed:      true,
		Sensitive:     true,
		Validators:    validators,
		PlanModifiers: modifiers,
		Default:       &nullDefault{},
	}
}

func SecretGenerated(optional bool, extras ...any) schema.StringAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.StringAttribute{
		Optional:      optional,
		Computed:      true,
		Sensitive:     true,
		Validators:    validators,
		PlanModifiers: append([]planmodifier.String{stringplanmodifier.UseStateForUnknown()}, modifiers...),
	}
}

func Optional(extras ...any) schema.StringAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.StringAttribute{
		Optional:      true,
		Computed:      true,
		Validators:    validators,
		PlanModifiers: append([]planmodifier.String{helpers.UseValidStateForUnknown()}, modifiers...),
	}
}

func Default(value string, extras ...any) schema.StringAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.StringAttribute{
		Optional:      true,
		Computed:      true,
		Validators:    validators,
		PlanModifiers: modifiers,
		Default:       stringdefault.StaticString(value),
	}
}

func Deprecated(message string, extras ...any) schema.StringAttribute {
	validators, modifiers := parseExtras(extras)
	return schema.StringAttribute{
		Optional:           true,
		Computed:           true,
		DeprecationMessage: message + " This attribute will be removed in a future version of the provider.",
		Validators:         validators,
		PlanModifiers:      modifiers,
		Default:            &nullDefault{},
	}
}

func Renamed(oldname, newname string, extras ...any) schema.StringAttribute {
	return Deprecated("The "+oldname+" attribute has been renamed, set the "+newname+" attribute instead.", extras...)
}

func Get(s Type, data map[string]any, key string) {
	if !s.IsNull() && !s.IsUnknown() {
		data[key] = s.ValueString()
	}
}

type SetOption int

const (
	SkipIfAlreadySet SetOption = iota
)

func Set(s *Type, data map[string]any, key string, options ...SetOption) {
	if v, ok := data[key].(string); ok {
		if s.ValueString() == "" || !slices.Contains(options, SkipIfAlreadySet) {
			*s = Value(v)
		}
	} else {
		Nil(s)
	}
}

func Nil(s *Type) {
	if s.IsUnknown() {
		*s = Value("")
	}
}

func parseExtras(extras []any) (validators []validator.String, modifiers []planmodifier.String) {
	for _, e := range extras {
		matched := false
		if validator, ok := e.(validator.String); ok {
			matched = true
			validators = append(validators, validator)
		}
		if modifier, ok := e.(planmodifier.String); ok {
			matched = true
			modifiers = append(modifiers, modifier)
		}
		if !matched {
			panic(fmt.Sprintf("unexpected extra value of type %T in string attribute", e))
		}
	}
	return
}
