package listtype

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.ListTypable                    = (*listNestedObjectTypeOf[struct{}])(nil)
	_ types.NestedObjectCollectionType         = (*listNestedObjectTypeOf[struct{}])(nil)
	_ basetypes.ListValuable                   = (*ListNestedObjectValueOf[struct{}])(nil)
	_ types.NestedObjectCollectionValue        = (*ListNestedObjectValueOf[struct{}])(nil)
	_ basetypes.ListValuableWithSemanticEquals = (*ListNestedObjectValueOf[struct{}])(nil)
)

type listSemanticEqualityFunc[T any] func(context.Context, ListNestedObjectValueOf[T], ListNestedObjectValueOf[T]) (bool, diag.Diagnostics)

// listNestedObjectTypeOf is the attribute type of a ListNestedObjectValueOf.
type listNestedObjectTypeOf[T any] struct {
	basetypes.ListType
	semanticEqualityFunc listSemanticEqualityFunc[T]
}

func NewListNestedObjectTypeOf[T any](ctx context.Context, f ...ListNestedObjectOfOption[T]) listNestedObjectTypeOf[T] {
	opts := newListNestedObjectOfOptions(f...)
	return listNestedObjectTypeOf[T]{
		ListType:             basetypes.ListType{ElemType: objtype.NewObjectTypeOfMust[T](ctx)},
		semanticEqualityFunc: opts.SemanticEqualityFunc,
	}
}

func (t listNestedObjectTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(listNestedObjectTypeOf[T])
	if !ok {
		return false
	}
	return t.ListType.Equal(other.ListType)
}

func (t listNestedObjectTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("ListNestedObjectTypeOf[%T]", zero)
}

func (t listNestedObjectTypeOf[T]) ValueFromList(ctx context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewListNestedObjectValueOfNull[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewListNestedObjectValueOfUnknown[T](ctx), diags
	}

	typ, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewListNestedObjectValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewListValue(typ, in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewListNestedObjectValueOfUnknown[T](ctx), diags
	}

	return ListNestedObjectValueOf[T]{
		ListValue:            v,
		semanticEqualityFunc: t.semanticEqualityFunc,
	}, diags
}

func (t listNestedObjectTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.ListType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	listValue, ok := attrValue.(basetypes.ListValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	listValuable, diags := t.ValueFromList(ctx, listValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ListValue to ListValuable: %v", diags)
	}

	return listValuable, nil
}

func (t listNestedObjectTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return ListNestedObjectValueOf[T]{semanticEqualityFunc: t.semanticEqualityFunc}
}

func (t listNestedObjectTypeOf[T]) NewObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return objtype.ObjectTypeNewObjectPtr[T](ctx)
}

func (t listNestedObjectTypeOf[T]) NewObjectSlice(ctx context.Context, len, cap int) (any, diag.Diagnostics) {
	return nestedObjectTypeNewObjectSlice[T](ctx, len, cap)
}

func (t listNestedObjectTypeOf[T]) NullValue(ctx context.Context) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	return NewListNestedObjectValueOfNull(ctx, WithSemanticEqualityFunc(t.semanticEqualityFunc)), diags
}

func (t listNestedObjectTypeOf[T]) ValueFromObjectPtr(ctx context.Context, ptr any) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := ptr.(*T); ok {
		v, d := newListNestedObjectValueOfPtr(ctx, v, t.semanticEqualityFunc)
		diags.Append(d...)
		return v, d
	}

	diags.Append(diag.NewErrorDiagnostic("Invalid pointer value", fmt.Sprintf("incorrect type: want %T, got %T", (*T)(nil), ptr)))
	return nil, diags
}

func (t listNestedObjectTypeOf[T]) ValueFromObjectSlice(ctx context.Context, slice any) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := slice.([]*T); ok {
		v, d := NewListNestedObjectValueOfSlice(ctx, v, t.semanticEqualityFunc)
		diags.Append(d...)
		return v, d
	}

	diags.Append(diag.NewErrorDiagnostic("Invalid slice value", fmt.Sprintf("incorrect type: want %T, got %T", (*[]T)(nil), slice)))
	return nil, diags
}

func nestedObjectTypeNewObjectSlice[T any](_ context.Context, len, cap int) ([]*T, diag.Diagnostics) { //nolint:unparam
	var diags diag.Diagnostics
	return make([]*T, len, cap), diags
}
