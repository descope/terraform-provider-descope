package emit

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

func ctyValue(ctx context.Context, value attr.Value) (cty.Value, error) {
	converted, err := value.ToTerraformValue(ctx)
	if err != nil {
		return cty.NilVal, err
	}
	return fromTerraform(converted)
}

func fromTerraform(value tftypes.Value) (cty.Value, error) {
	typ := value.Type()
	if value.IsNull() {
		converted, err := fromTerraformType(typ)
		if err != nil {
			return cty.NilVal, err
		}
		return cty.NullVal(converted), nil
	}

	switch {
	case typ.Is(tftypes.String):
		var s string
		if err := value.As(&s); err != nil {
			return cty.NilVal, err
		}
		return cty.StringVal(s), nil
	case typ.Is(tftypes.Bool):
		var b bool
		if err := value.As(&b); err != nil {
			return cty.NilVal, err
		}
		return cty.BoolVal(b), nil
	case typ.Is(tftypes.Number):
		n := new(big.Float)
		if err := value.As(&n); err != nil {
			return cty.NilVal, err
		}
		return cty.NumberVal(n), nil
	case typ.Is(tftypes.List{}) || typ.Is(tftypes.Set{}) || typ.Is(tftypes.Tuple{}):
		var values []tftypes.Value
		if err := value.As(&values); err != nil {
			return cty.NilVal, err
		}
		elements := make([]cty.Value, 0, len(values))
		for _, v := range values {
			element, err := fromTerraform(v)
			if err != nil {
				return cty.NilVal, err
			}
			elements = append(elements, element)
		}
		switch t := typ.(type) {
		case tftypes.List:
			if len(elements) == 0 {
				elementType, err := fromTerraformType(t.ElementType)
				if err != nil {
					return cty.NilVal, err
				}
				return cty.ListValEmpty(elementType), nil
			}
			return cty.ListVal(elements), nil
		case tftypes.Set:
			if len(elements) == 0 {
				elementType, err := fromTerraformType(t.ElementType)
				if err != nil {
					return cty.NilVal, err
				}
				return cty.SetValEmpty(elementType), nil
			}
			return cty.SetVal(elements), nil
		default:
			return cty.TupleVal(elements), nil
		}
	case typ.Is(tftypes.Map{}):
		var values map[string]tftypes.Value
		if err := value.As(&values); err != nil {
			return cty.NilVal, err
		}
		if len(values) == 0 {
			mapType, ok := typ.(tftypes.Map)
			if !ok {
				return cty.NilVal, fmt.Errorf("unexpected map type %s", typ)
			}
			elementType, err := fromTerraformType(mapType.ElementType)
			if err != nil {
				return cty.NilVal, err
			}
			return cty.MapValEmpty(elementType), nil
		}
		elements := map[string]cty.Value{}
		for key, v := range values {
			element, err := fromTerraform(v)
			if err != nil {
				return cty.NilVal, err
			}
			elements[key] = element
		}
		return cty.MapVal(elements), nil
	case typ.Is(tftypes.Object{}):
		var values map[string]tftypes.Value
		if err := value.As(&values); err != nil {
			return cty.NilVal, err
		}
		elements := map[string]cty.Value{}
		for key, v := range values {
			element, err := fromTerraform(v)
			if err != nil {
				return cty.NilVal, err
			}
			elements[key] = element
		}
		return cty.ObjectVal(elements), nil
	}
	return cty.NilVal, fmt.Errorf("unsupported value type %s", typ)
}

func fromTerraformType(typ tftypes.Type) (cty.Type, error) {
	switch {
	case typ.Is(tftypes.String):
		return cty.String, nil
	case typ.Is(tftypes.Bool):
		return cty.Bool, nil
	case typ.Is(tftypes.Number):
		return cty.Number, nil
	}
	switch t := typ.(type) {
	case tftypes.List:
		element, err := fromTerraformType(t.ElementType)
		if err != nil {
			return cty.NilType, err
		}
		return cty.List(element), nil
	case tftypes.Set:
		element, err := fromTerraformType(t.ElementType)
		if err != nil {
			return cty.NilType, err
		}
		return cty.Set(element), nil
	case tftypes.Map:
		element, err := fromTerraformType(t.ElementType)
		if err != nil {
			return cty.NilType, err
		}
		return cty.Map(element), nil
	case tftypes.Object:
		attributeTypes := map[string]cty.Type{}
		for name, attributeType := range t.AttributeTypes {
			converted, err := fromTerraformType(attributeType)
			if err != nil {
				return cty.NilType, err
			}
			attributeTypes[name] = converted
		}
		return cty.Object(attributeTypes), nil
	}
	return cty.NilType, fmt.Errorf("unsupported type %s", typ)
}

func ValueJSON(ctx context.Context, value attr.Value) ([]byte, error) {
	converted, err := ctyValue(ctx, value)
	if err != nil {
		return nil, err
	}
	return ctyjson.Marshal(converted, converted.Type())
}
