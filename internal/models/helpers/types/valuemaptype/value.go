// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package valuemaptype

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Value            = (*MapValueOf[basetypes.StringValue])(nil)
	_ basetypes.MapValuable = (*MapValueOf[basetypes.StringValue])(nil)
)

type MapValueOf[T attr.Value] struct {
	basetypes.MapValue
}

func (v MapValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(MapValueOf[T])
	if !ok {
		return false
	}
	return v.MapValue.Equal(other.MapValue)
}

func (v MapValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewMapTypeOf[T](ctx)
}

func (v MapValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.MapValue.ToTerraformValue(ctx)
}

func (v MapValueOf[T]) ToMap(ctx context.Context) (map[string]T, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := map[string]T{}
	d := v.ElementsAs(ctx, &values, false)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return values, diags
}

func (v MapValueOf[T]) ToMapMust(ctx context.Context) map[string]T {
	return helpers.Must(v.ToMap(ctx))
}

func (v MapValueOf[T]) IsEmpty() bool {
	return len(v.MapValue.Elements()) == 0
}

func NewMapValueOfNull[T attr.Value](ctx context.Context) MapValueOf[T] {
	return MapValueOf[T]{MapValue: basetypes.NewMapNull(types.NewAttrTypeOf[T](ctx))}
}

func NewMapValueOfUnknown[T attr.Value](ctx context.Context) MapValueOf[T] {
	return MapValueOf[T]{MapValue: basetypes.NewMapUnknown(types.NewAttrTypeOf[T](ctx))}
}

func NewMapValueOf[T attr.Value](ctx context.Context, elements map[string]attr.Value) (MapValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	v, d := basetypes.NewMapValue(types.NewAttrTypeOf[T](ctx), elements)
	diags.Append(d...)
	if diags.HasError() {
		return NewMapValueOfUnknown[T](ctx), diags
	}

	return MapValueOf[T]{MapValue: v}, diags
}

func NewMapValueOfMust[T attr.Value](ctx context.Context, elements map[string]attr.Value) MapValueOf[T] {
	return helpers.Must(NewMapValueOf[T](ctx, elements))
}
