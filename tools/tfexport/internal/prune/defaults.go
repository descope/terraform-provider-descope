package prune

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func AttributeDefault(ctx context.Context, attribute schema.Attribute) (attr.Value, bool) {
	switch a := attribute.(type) {
	case schema.StringAttribute:
		if a.Default != nil {
			resp := defaults.StringResponse{}
			a.Default.DefaultString(ctx, defaults.StringRequest{}, &resp)
			return resp.PlanValue, true
		}
	case schema.BoolAttribute:
		if a.Default != nil {
			resp := defaults.BoolResponse{}
			a.Default.DefaultBool(ctx, defaults.BoolRequest{}, &resp)
			return resp.PlanValue, true
		}
	case schema.Int64Attribute:
		if a.Default != nil {
			resp := defaults.Int64Response{}
			a.Default.DefaultInt64(ctx, defaults.Int64Request{}, &resp)
			return resp.PlanValue, true
		}
	case schema.Float64Attribute:
		if a.Default != nil {
			resp := defaults.Float64Response{}
			a.Default.DefaultFloat64(ctx, defaults.Float64Request{}, &resp)
			return resp.PlanValue, true
		}
	case schema.ListAttribute:
		if a.Default != nil {
			resp := defaults.ListResponse{}
			a.Default.DefaultList(ctx, defaults.ListRequest{}, &resp)
			return resp.PlanValue, true
		}
	case schema.SetAttribute:
		if a.Default != nil {
			resp := defaults.SetResponse{}
			a.Default.DefaultSet(ctx, defaults.SetRequest{}, &resp)
			return resp.PlanValue, true
		}
	case schema.MapAttribute:
		if a.Default != nil {
			resp := defaults.MapResponse{}
			a.Default.DefaultMap(ctx, defaults.MapRequest{}, &resp)
			return resp.PlanValue, true
		}
	case schema.ListNestedAttribute:
		if a.Default != nil {
			resp := defaults.ListResponse{}
			a.Default.DefaultList(ctx, defaults.ListRequest{}, &resp)
			return resp.PlanValue, true
		}
	case schema.SetNestedAttribute:
		if a.Default != nil {
			resp := defaults.SetResponse{}
			a.Default.DefaultSet(ctx, defaults.SetRequest{}, &resp)
			return resp.PlanValue, true
		}
	case schema.MapNestedAttribute:
		if a.Default != nil {
			resp := defaults.MapResponse{}
			a.Default.DefaultMap(ctx, defaults.MapRequest{}, &resp)
			return resp.PlanValue, true
		}
	case schema.SingleNestedAttribute:
		if a.Default != nil {
			resp := defaults.ObjectResponse{}
			a.Default.DefaultObject(ctx, defaults.ObjectRequest{}, &resp)
			return resp.PlanValue, true
		}
	}
	return nil, false
}

func isServerZeroValue(ctx context.Context, attribute schema.Attribute, value attr.Value) bool {
	if attribute.IsRequired() || !attribute.IsComputed() {
		return false
	}
	raw, err := value.ToTerraformValue(ctx)
	if err != nil {
		return false
	}
	switch {
	case raw.Type().Is(tftypes.String):
		var s string
		return raw.As(&s) == nil && s == ""
	case raw.Type().Is(tftypes.List{}), raw.Type().Is(tftypes.Set{}), raw.Type().Is(tftypes.Map{}):
		var elements []tftypes.Value
		if raw.As(&elements) == nil {
			return len(elements) == 0
		}
		var entries map[string]tftypes.Value
		return raw.As(&entries) == nil && len(entries) == 0
	}
	return false
}

func valuesEqual(ctx context.Context, a, b attr.Value) bool {
	av, err := a.ToTerraformValue(ctx)
	if err != nil {
		return false
	}
	bv, err := b.ToTerraformValue(ctx)
	if err != nil {
		return false
	}
	return av.Equal(bv)
}
