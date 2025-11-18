package durationattr

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type = types.String

func Value(value string) Type {
	return types.StringValue(value)
}

func Required(validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Required:   true,
		Validators: append([]validator.String{formatValidator}, validators...),
	}
}

func Optional(validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:      true,
		Computed:      true,
		Validators:    append([]validator.String{formatValidator}, validators...),
		PlanModifiers: []planmodifier.String{helpers.UseValidStateForUnknown()},
	}
}

func Default(value string, validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:   true,
		Computed:   true,
		Validators: append([]validator.String{formatValidator}, validators...),
		Default:    stringdefault.StaticString(value),
	}
}

func Get(s types.String, data map[string]any, key string) {
	if !s.IsNull() && !s.IsUnknown() {
		num, unit, _ := parseString(s.ValueString())
		data[key] = num
		data[key+"Unit"] = unit
	}
}

func Set(s *types.String, data map[string]any, key string) {
	num, hasNum := data[key].(int64)
	if v, ok := data[key].(float64); ok {
		hasNum = true
		num = int64(v)
	}
	unit, hasUnit := data[key+"Unit"].(string)
	if !hasNum || !hasUnit {
		return
	}
	value := composeString(num, unit)
	if value != s.ValueString()+"s" { // don't overwrite singular with plural
		*s = types.StringValue(value)
	}
}

// Utils

var units = []string{"seconds", "minutes", "hours", "days", "weeks"}

func composeString(num int64, unit string) string {
	return fmt.Sprintf("%d %s", num, unit)
}

func parseString(s string) (num int64, unit string, ok bool) {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return
	}
	for _, r := range parts[0] {
		if r < '0' || r > '9' {
			return
		}
	}
	num, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil || num > 1000 {
		return
	}
	unit = strings.TrimSuffix(parts[1], "s") + "s"
	if !slices.Contains(units, unit) {
		return
	}
	return num, unit, true
}
