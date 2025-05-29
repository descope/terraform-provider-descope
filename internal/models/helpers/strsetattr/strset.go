package strsetattr

import (
	"context"
	"fmt"
	"iter"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/valuesettype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type = valuesettype.SetValueOf[types.String]

func Value(value []string) Type {
	return convertStringSliceToTerraformValue(context.Background(), value)
}

func Required(validators ...validator.Set) schema.SetAttribute {
	return schema.SetAttribute{
		Required:    true,
		CustomType:  valuesettype.StringSetType,
		ElementType: types.StringType,
		Validators:  validators,
	}
}

func Optional(validators ...validator.Set) schema.SetAttribute {
	return schema.SetAttribute{
		Optional:      true,
		Computed:      true,
		CustomType:    valuesettype.StringSetType,
		ElementType:   types.StringType,
		Validators:    validators,
		PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
	}
}

func Default(validators ...validator.Set) schema.SetAttribute {
	return schema.SetAttribute{
		Optional:    true,
		Computed:    true,
		CustomType:  valuesettype.StringSetType,
		ElementType: types.StringType,
		Validators:  validators,
		Default:     setdefault.StaticValue(Value([]string{}).SetValue),
	}
}

func Get(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsUnknown() {
		return
	}

	values := s.ToSliceMust(h.Ctx)
	data[key] = helpers.ConvertTerraformSliceToStringSlice(values)
}

func Set(s *Type, data map[string]any, key string, h *helpers.Handler) {
	values := helpers.GetStringSlice(data, key)
	*s = convertStringSliceToTerraformValue(h.Ctx, values)
}

func GetCommaSeparated(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsUnknown() {
		return
	}

	values := s.ToSliceMust(h.Ctx)
	data[key] = strings.Join(helpers.ConvertTerraformSliceToStringSlice(values), ",")
}

func SetCommaSeparated(s *Type, data map[string]any, key string, h *helpers.Handler) {
	values := helpers.GetCommaSeparatedStringSlice(data, key)
	*s = convertStringSliceToTerraformValue(h.Ctx, values)
}

func Iterator(l Type, h *helpers.Handler) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range l.Elements() {
			if v.IsNull() || v.IsUnknown() {
				continue
			}

			s, ok := v.(types.String)
			if !ok {
				h.Diagnostics.Append(diag.NewErrorDiagnostic("Unexpected Value Type", fmt.Sprintf("Expected string type, found %T", v)))
				continue
			}

			if !yield(s.ValueString()) {
				break
			}
		}
	}
}

func convertStringSliceToTerraformValue(ctx context.Context, values []string) Type {
	var elements []attr.Value
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	return valuesettype.NewSetValueOfMust[types.String](ctx, elements)
}
