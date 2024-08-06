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
