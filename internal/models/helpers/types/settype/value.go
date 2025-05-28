package settype

import (
	"context"
	"fmt"
	"iter"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Value                              = (*SetNestedObjectValueOf[struct{}])(nil)
	_ basetypes.SetValuable                   = (*SetNestedObjectValueOf[struct{}])(nil)
	_ basetypes.SetValuableWithSemanticEquals = (*SetNestedObjectValueOf[struct{}])(nil)
	_ types.NestedObjectValue                 = (*SetNestedObjectValueOf[struct{}])(nil)
	_ types.NestedObjectCollectionValue       = (*SetNestedObjectValueOf[struct{}])(nil)
)

type SetNestedObjectValueOf[T any] struct {
	basetypes.SetValue
	semanticEqualityFunc setSemanticEqualityFunc[T]
}

func (v SetNestedObjectValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(SetNestedObjectValueOf[T])
	if !ok {
		return false
	}
	return v.SetValue.Equal(other.SetValue)
}

func (v SetNestedObjectValueOf[T]) SetSemanticEquals(ctx context.Context, newValuable basetypes.SetValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.semanticEqualityFunc == nil {
		return false, diags
	}

	newValue, ok := newValuable.(SetNestedObjectValueOf[T])
	if !ok {
		diags.AddError("SetSemanticEquals", fmt.Sprintf("unexpected value type of %T", newValuable))
		return false, diags
	}

	return v.semanticEqualityFunc(ctx, v, newValue)
}

func (v SetNestedObjectValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewSetNestedObjectTypeOfMust[T](ctx)
}

func (v SetNestedObjectValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.SetValue.ToTerraformValue(ctx)
}

func (v SetNestedObjectValueOf[T]) ToObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return v.ToPtr(ctx)
}

func (v SetNestedObjectValueOf[T]) ToObjectSlice(ctx context.Context) (any, diag.Diagnostics) {
	return v.ToSlice(ctx)
}

func (v SetNestedObjectValueOf[T]) ToPtr(ctx context.Context) (*T, diag.Diagnostics) {
	return nestedObjectValueObjectPtr[T](ctx, v.SetValue)
}

func (v SetNestedObjectValueOf[T]) ToSlice(ctx context.Context) ([]*T, diag.Diagnostics) {
	return nestedObjectValueObjectSlice[T](ctx, v.SetValue)
}

func (v SetNestedObjectValueOf[T]) IsEmpty() bool {
	return len(v.SetValue.Elements()) == 0
}

func (v SetNestedObjectValueOf[T]) ImmutableIterator(ctx context.Context) iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for _, v := range v.Elements() {
			ptr, diags := objtype.ObjectValueObjectPtr[T](ctx, v)
			if diags.HasError() {
				continue
			}
			if !yield(ptr) {
				break
			}
		}
	}
}

func (v *SetNestedObjectValueOf[T]) MutableIterator(ctx context.Context) iter.Seq[*T] {
	return func(yield func(*T) bool) {
		s, _ := v.ToSlice(ctx)
		if s == nil {
			s = []*T{}
		}

		for _, v := range s {
			if !yield(v) {
				break
			}
		}

		*v, _ = newSetNestedObjectValueOf(ctx, s, v.semanticEqualityFunc)
	}
}

func nestedObjectValueObjectPtr[T any](ctx context.Context, val types.ValueWithElements) (*T, diag.Diagnostics) { // TODO
	var diags diag.Diagnostics

	elements := val.Elements()
	switch n := len(elements); n {
	case 0:
		return nil, diags
	case 1:
		ptr, d := objtype.ObjectValueObjectPtr[T](ctx, elements[0])
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		return ptr, diags
	default:
		diags.Append(diag.NewErrorDiagnostic("Invalid list/set", fmt.Sprintf("too many elements: want 1, got %d", n)))
		return nil, diags
	}
}

func nestedObjectValueObjectSlice[T any](ctx context.Context, val types.ValueWithElements) ([]*T, diag.Diagnostics) {
	var diags diag.Diagnostics

	elements := val.Elements()
	n := len(elements)
	slice := make([]*T, n)
	for i := range n {
		ptr, d := objtype.ObjectValueObjectPtr[T](ctx, elements[i])
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		slice[i] = ptr
	}

	return slice, diags
}

func NewSetNestedObjectValueOfNull[T any](ctx context.Context, f ...SetNestedObjectOfOption[T]) SetNestedObjectValueOf[T] {
	opts := newSetNestedObjectOfOptions(f...)
	return SetNestedObjectValueOf[T]{SetValue: basetypes.NewSetNull(objtype.NewObjectTypeOfMust[T](ctx)), semanticEqualityFunc: opts.SemanticEqualityFunc}
}

func NewSetNestedObjectValueOfUnknown[T any](ctx context.Context) SetNestedObjectValueOf[T] {
	return SetNestedObjectValueOf[T]{SetValue: basetypes.NewSetUnknown(objtype.NewObjectTypeOfMust[T](ctx))}
}

func NewSetNestedObjectValueOfPtr[T any](ctx context.Context, t *T, f ...SetNestedObjectOfOption[T]) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	opts := newSetNestedObjectOfOptions(f...)
	return newSetNestedObjectValueOfPtr(ctx, t, opts.SemanticEqualityFunc)
}

func newSetNestedObjectValueOfPtr[T any](ctx context.Context, t *T, f setSemanticEqualityFunc[T]) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	return NewSetNestedObjectValueOfSlice(ctx, []*T{t}, f)
}

func NewSetNestedObjectValueOfPtrMust[T any](ctx context.Context, t *T, f ...SetNestedObjectOfOption[T]) SetNestedObjectValueOf[T] {
	opts := newSetNestedObjectOfOptions(f...)
	return types.Must(newSetNestedObjectValueOfPtr(ctx, t, opts.SemanticEqualityFunc))
}

func NewSetNestedObjectValueOfSlice[T any](ctx context.Context, ts []*T, f setSemanticEqualityFunc[T]) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	return newSetNestedObjectValueOf(ctx, ts, f)
}

func NewSetNestedObjectValueOfSliceMust[T any](ctx context.Context, ts []*T, f ...SetNestedObjectOfOption[T]) SetNestedObjectValueOf[T] {
	opts := newSetNestedObjectOfOptions(f...)
	return types.Must(NewSetNestedObjectValueOfSlice(ctx, ts, opts.SemanticEqualityFunc))
}

func newSetNestedObjectValueOf[T any](ctx context.Context, elements []*T, f setSemanticEqualityFunc[T]) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	typ, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewSetNestedObjectValueOfUnknown[T](ctx), diags
	}

	values := []attr.Value{}
	for _, v := range elements {
		values = append(values, objtype.NewObjectValueOfMust(ctx, v))
	}

	v, d := basetypes.NewSetValue(typ, values)
	diags.Append(d...)
	if diags.HasError() {
		return NewSetNestedObjectValueOfUnknown[T](ctx), diags
	}

	return SetNestedObjectValueOf[T]{SetValue: v, semanticEqualityFunc: f}, diags
}
