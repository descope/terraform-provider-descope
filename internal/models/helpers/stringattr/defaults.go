package stringattr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NullDefault() defaults.String {
	return &nullDefault{}
}

type nullDefault struct {
}

func (d nullDefault) Description(_ context.Context) string {
	return "value defaults to null"
}

func (d nullDefault) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d nullDefault) DefaultString(_ context.Context, _ defaults.StringRequest, resp *defaults.StringResponse) {
	resp.PlanValue = types.StringNull()
}

// TODO
// func DefaultModifier(value string) planmodifier.String {
// 	return defaultModifier{value: value}
// }

// type defaultModifier struct {
// 	value string
// }

// func (d defaultModifier) Description(_ context.Context) string {
// 	return "sets the default value for a string attribute if it is unknown in the plan"
// }

// func (d defaultModifier) MarkdownDescription(ctx context.Context) string {
// 	return d.Description(ctx)
// }

// func (d defaultModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
// 	if req.PlanValue.IsUnknown() {
// 		resp.PlanValue = Value(d.value)
// 	}
// }
