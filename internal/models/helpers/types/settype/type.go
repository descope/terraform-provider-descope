package settype

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
	_ attr.Type                        = (*setNestedObjectTypeOf[struct{}])(nil)
	_ attr.TypeWithElementType         = (*setNestedObjectTypeOf[struct{}])(nil)
	_ basetypes.SetTypable             = (*setNestedObjectTypeOf[struct{}])(nil)
	_ types.NestedObjectType           = (*setNestedObjectTypeOf[struct{}])(nil)
	_ types.NestedObjectCollectionType = (*setNestedObjectTypeOf[struct{}])(nil)
)

type setSemanticEqualityFunc[T any] func(context.Context, SetNestedObjectValueOf[T], SetNestedObjectValueOf[T]) (bool, diag.Diagnostics)

type setNestedObjectTypeOf[T any] struct {
	basetypes.SetType
	semanticEqualityFunc setSemanticEqualityFunc[T]
}

func NewSetNestedObjectTypeOf[T any](ctx context.Context, f ...SetNestedObjectOfOption[T]) (setNestedObjectTypeOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	elemType, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return setNestedObjectTypeOf[T]{}, diags
	}

	opts := newSetNestedObjectOfOptions(f...)

	return setNestedObjectTypeOf[T]{
		SetType:              basetypes.SetType{ElemType: elemType},
		semanticEqualityFunc: opts.SemanticEqualityFunc,
	}, diags
}

func NewSetNestedObjectTypeOfMust[T any](ctx context.Context, f ...SetNestedObjectOfOption[T]) setNestedObjectTypeOf[T] {
	return types.Must(NewSetNestedObjectTypeOf(ctx, f...))
}

func (t setNestedObjectTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(setNestedObjectTypeOf[T])
	if !ok {
		return false
	}
	return t.SetType.Equal(other.SetType)
}

func (t setNestedObjectTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("SetNestedObjectTypeOf[%T]", zero)
}

func (t setNestedObjectTypeOf[T]) ValueFromSet(ctx context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewSetNestedObjectValueOfNull[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewSetNestedObjectValueOfUnknown[T](ctx), diags
	}

	typ, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewSetNestedObjectValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewSetValue(typ, in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewSetNestedObjectValueOfUnknown[T](ctx), diags
	}

	return SetNestedObjectValueOf[T]{
		SetValue:             v,
		semanticEqualityFunc: t.semanticEqualityFunc,
	}, diags
}

func (t setNestedObjectTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.SetType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	setValue, ok := attrValue.(basetypes.SetValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	setValuable, diags := t.ValueFromSet(ctx, setValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting SetValue to SetValuable: %v", diags)
	}

	return setValuable, nil
}

func (t setNestedObjectTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return SetNestedObjectValueOf[T]{semanticEqualityFunc: t.semanticEqualityFunc}
}

func (t setNestedObjectTypeOf[T]) NewObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return objtype.ObjectTypeNewObjectPtr[T](ctx)
}

func (t setNestedObjectTypeOf[T]) NewObjectSlice(ctx context.Context, len, cap int) (any, diag.Diagnostics) {
	return nestedObjectTypeNewObjectSlice[T](ctx, len, cap)
}

func (t setNestedObjectTypeOf[T]) NullValue(ctx context.Context) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	return NewSetNestedObjectValueOfNull(ctx, WithSemanticEqualityFunc(t.semanticEqualityFunc)), diags
}

func (t setNestedObjectTypeOf[T]) ValueFromObjectPtr(ctx context.Context, ptr any) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := ptr.(*T); ok {
		v, d := newSetNestedObjectValueOfPtr(ctx, v, t.semanticEqualityFunc)
		diags.Append(d...)
		return v, d
	}

	diags.Append(diag.NewErrorDiagnostic("Invalid pointer value", fmt.Sprintf("incorrect type: want %T, got %T", (*T)(nil), ptr)))
	return nil, diags
}

func (t setNestedObjectTypeOf[T]) ValueFromObjectSlice(ctx context.Context, slice any) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := slice.([]*T); ok {
		v, d := NewSetNestedObjectValueOfSlice(ctx, v, t.semanticEqualityFunc)
		diags.Append(d...)
		return v, d
	}

	diags.Append(diag.NewErrorDiagnostic("Invalid slice value", fmt.Sprintf("incorrect type: want %T, got %T", (*[]T)(nil), slice)))
	return nil, diags
}

func nestedObjectTypeNewObjectSlice[T any](_ context.Context, len, cap int) ([]*T, diag.Diagnostics) {
	var diags diag.Diagnostics
	return make([]*T, len, cap), diags
}
