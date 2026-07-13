package boolattr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type nullDefault struct {
}

func (d nullDefault) Description(_ context.Context) string {
	return "value defaults to null"
}

func (d nullDefault) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d nullDefault) DefaultBool(_ context.Context, _ defaults.BoolRequest, resp *defaults.BoolResponse) {
	resp.PlanValue = types.BoolNull()
}
