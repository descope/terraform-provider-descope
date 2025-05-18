package listtype

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ListNestedObjectValueOf represents a Terraform Plugin Framework List value whose elements are of type `ObjectTypeOf[T]`.
type ListNestedObjectValueOf[T any] struct {
	basetypes.ListValue
	semanticEqualityFunc listSemanticEqualityFunc[T]
}

func (v ListNestedObjectValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(ListNestedObjectValueOf[T])
	if !ok {
		return false
	}
	return v.ListValue.Equal(other.ListValue)
}

func (v ListNestedObjectValueOf[T]) ListSemanticEquals(ctx context.Context, newValuable basetypes.ListValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	// returning false will fall back to framework defined semantic equality checks
	if v.semanticEqualityFunc == nil {
		return false, diags
	}

	newValue, ok := newValuable.(ListNestedObjectValueOf[T])
	if !ok {
		diags.AddError("ListSemanticEquals", fmt.Sprintf("unexpected value type of %T", newValuable))
		return false, diags
	}

	return v.semanticEqualityFunc(ctx, v, newValue)
}

func (v ListNestedObjectValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewListNestedObjectTypeOf[T](ctx)
}

func (v ListNestedObjectValueOf[T]) ToObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return v.ToPtr(ctx)
}

func (v ListNestedObjectValueOf[T]) ToObjectSlice(ctx context.Context) (any, diag.Diagnostics) {
	return v.ToSlice(ctx)
}

// ToPtr returns a pointer to the single element of a ListNestedObject.
func (v ListNestedObjectValueOf[T]) ToPtr(ctx context.Context) (*T, diag.Diagnostics) {
	return nestedObjectValueObjectPtr[T](ctx, v.ListValue)
}

// ToSlice returns a slice of pointers to the elements of a ListNestedObject.
func (v ListNestedObjectValueOf[T]) ToSlice(ctx context.Context) ([]*T, diag.Diagnostics) {
	return nestedObjectValueObjectSlice[T](ctx, v.ListValue)
}

func nestedObjectValueObjectPtr[T any](ctx context.Context, val types.ValueWithElements) (*T, diag.Diagnostics) {
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

func NewListNestedObjectValueOfNull[T any](ctx context.Context, f ...ListNestedObjectOfOption[T]) ListNestedObjectValueOf[T] {
	opts := newListNestedObjectOfOptions(f...)
	return ListNestedObjectValueOf[T]{ListValue: basetypes.NewListNull(objtype.NewObjectTypeOfMust[T](ctx)), semanticEqualityFunc: opts.SemanticEqualityFunc}
}

func NewListNestedObjectValueOfUnknown[T any](ctx context.Context) ListNestedObjectValueOf[T] {
	return ListNestedObjectValueOf[T]{ListValue: basetypes.NewListUnknown(objtype.NewObjectTypeOfMust[T](ctx))}
}

func NewListNestedObjectValueOfPtr[T any](ctx context.Context, t *T, f ...ListNestedObjectOfOption[T]) (ListNestedObjectValueOf[T], diag.Diagnostics) {
	opts := newListNestedObjectOfOptions(f...)
	return newListNestedObjectValueOfPtr(ctx, t, opts.SemanticEqualityFunc)
}

func newListNestedObjectValueOfPtr[T any](ctx context.Context, t *T, f listSemanticEqualityFunc[T]) (ListNestedObjectValueOf[T], diag.Diagnostics) {
	return NewListNestedObjectValueOfSlice(ctx, []*T{t}, f)
}

func NewListNestedObjectValueOfPtrMust[T any](ctx context.Context, t *T, f ...ListNestedObjectOfOption[T]) ListNestedObjectValueOf[T] {
	opts := newListNestedObjectOfOptions(f...)
	return types.Must(newListNestedObjectValueOfPtr(ctx, t, opts.SemanticEqualityFunc))
}

func NewListNestedObjectValueOfSlice[T any](ctx context.Context, ts []*T, f listSemanticEqualityFunc[T]) (ListNestedObjectValueOf[T], diag.Diagnostics) {
	return newListNestedObjectValueOf(ctx, ts, f)
}

func NewListNestedObjectValueOfSliceMust[T any](ctx context.Context, ts []*T, f ...ListNestedObjectOfOption[T]) ListNestedObjectValueOf[T] {
	opts := newListNestedObjectOfOptions(f...)
	return types.Must(NewListNestedObjectValueOfSlice(ctx, ts, opts.SemanticEqualityFunc))
}

func NewListNestedObjectValueOfValueSlice[T any](ctx context.Context, ts []T, f ...ListNestedObjectOfOption[T]) (ListNestedObjectValueOf[T], diag.Diagnostics) {
	opts := newListNestedObjectOfOptions(f...)
	return newListNestedObjectValueOf[T](ctx, ts, opts.SemanticEqualityFunc)
}

func NewListNestedObjectValueOfValueSliceMust[T any](ctx context.Context, ts []T, f ...ListNestedObjectOfOption[T]) ListNestedObjectValueOf[T] {
	return types.Must(NewListNestedObjectValueOfValueSlice(ctx, ts, f...))
}

func newListNestedObjectValueOf[T any](ctx context.Context, elements any, f listSemanticEqualityFunc[T]) (ListNestedObjectValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	typ, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewListNestedObjectValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewListValueFrom(ctx, typ, elements)
	diags.Append(d...)
	if diags.HasError() {
		return NewListNestedObjectValueOfUnknown[T](ctx), diags
	}

	return ListNestedObjectValueOf[T]{ListValue: v, semanticEqualityFunc: f}, diags
}
