package maptype

import (
	"context"
	"iter"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Value            = (*MapNestedObjectValueOf[struct{}])(nil)
	_ basetypes.MapValuable = (*MapNestedObjectValueOf[struct{}])(nil)
)

type MapNestedObjectValueOf[T any] struct {
	basetypes.MapValue
}

func (v MapNestedObjectValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(MapNestedObjectValueOf[T])
	if !ok {
		return false
	}
	return v.MapValue.Equal(other.MapValue)
}

func (v MapNestedObjectValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewMapNestedObjectTypeOfMust[T](ctx)
}

func (v MapNestedObjectValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.MapValue.ToTerraformValue(ctx)
}

func (v MapNestedObjectValueOf[T]) ToObjectMap(ctx context.Context) (any, diag.Diagnostics) {
	return v.ToMap(ctx)
}

func (v MapNestedObjectValueOf[T]) ToMap(ctx context.Context) (map[string]*T, diag.Diagnostics) {
	return nestedObjectValueObjectMap(ctx, v)
}

func (v MapNestedObjectValueOf[T]) ToMapMust(ctx context.Context) map[string]*T {
	return types.Must(nestedObjectValueObjectMap(ctx, v))
}

func (v MapNestedObjectValueOf[T]) IsEmpty() bool {
	return len(v.MapValue.Elements()) == 0
}

func (v MapNestedObjectValueOf[T]) ImmutableIterator(ctx context.Context) iter.Seq2[string, *T] {
	return func(yield func(string, *T) bool) {
		for k, v := range v.Elements() {
			ptr, diags := objtype.ObjectValueObjectPtr[T](ctx, v)
			if diags.HasError() {
				continue
			}
			if !yield(k, ptr) {
				break
			}
		}
	}
}

func (v *MapNestedObjectValueOf[T]) MutableIterator(ctx context.Context) iter.Seq2[string, *T] {
	return func(yield func(string, *T) bool) {
		m, _ := v.ToMap(ctx)
		if m == nil {
			m = map[string]*T{}
		}

		for k, v := range m {
			if !yield(k, v) {
				break
			}
		}

		*v, _ = newMapNestedObjectValueOf(ctx, m)
	}
}

func nestedObjectValueObjectMap[T any](ctx context.Context, val MapNestedObjectValueOf[T]) (map[string]*T, diag.Diagnostics) {
	var diags diag.Diagnostics

	elements := val.Elements()
	result := make(map[string]*T, len(elements))
	for k, v := range elements {
		ptr, d := objtype.ObjectValueObjectPtr[T](ctx, v)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		result[k] = ptr
	}

	return result, diags
}

func NewMapNestedObjectValueOfNull[T any](ctx context.Context) MapNestedObjectValueOf[T] {
	return MapNestedObjectValueOf[T]{MapValue: basetypes.NewMapNull(objtype.NewObjectTypeOfMust[T](ctx))}
}

func NewMapNestedObjectValueOfUnknown[T any](ctx context.Context) MapNestedObjectValueOf[T] {
	return MapNestedObjectValueOf[T]{MapValue: basetypes.NewMapUnknown(objtype.NewObjectTypeOfMust[T](ctx))}
}

func NewMapNestedObjectValueOfMap[T any](ctx context.Context, ts map[string]*T) (MapNestedObjectValueOf[T], diag.Diagnostics) {
	return newMapNestedObjectValueOf(ctx, ts)
}

func NewMapNestedObjectValueOfMapMust[T any](ctx context.Context, ts map[string]*T) MapNestedObjectValueOf[T] {
	return types.Must(NewMapNestedObjectValueOfMap(ctx, ts))
}

func newMapNestedObjectValueOf[T any](ctx context.Context, elements map[string]*T) (MapNestedObjectValueOf[T], diag.Diagnostics) {
	values := map[string]attr.Value{}
	for k, v := range elements {
		values[k] = objtype.NewObjectValueOfMust(ctx, v)
	}
	return ValueOf[T](ctx, values)
}

func ValueOf[T any](ctx context.Context, elements map[string]attr.Value) (MapNestedObjectValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	typ, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewMapNestedObjectValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewMapValue(typ, elements)
	diags.Append(d...)
	if diags.HasError() {
		return NewMapNestedObjectValueOfUnknown[T](ctx), diags
	}

	return MapNestedObjectValueOf[T]{MapValue: v}, diags

}
