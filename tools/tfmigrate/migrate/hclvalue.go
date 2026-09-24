package migrate

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/zclconf/go-cty/cty"
)

// expr is a generated configuration value. Most values are plain literals, but a value holding a secret or an
// external payload has to be an expression instead, and a nested object can mix the two.
type expr struct {
	lit    cty.Value
	raw    string
	obj    map[string]*expr
	list   []*expr
	isObj  bool
	isList bool
}

func literalExpr(value cty.Value) *expr { return &expr{lit: value} }

func rawExpr(text string) *expr { return &expr{raw: text} }

func objectExpr(attrs map[string]*expr) *expr { return &expr{obj: attrs, isObj: true} }

func listExpr(items []*expr) *expr { return &expr{list: items, isList: true} }

// literal reports whether the expression is a plain value that can be written with SetAttributeValue.
func (e *expr) literal() (cty.Value, bool) {
	if e.isObj || e.isList || e.raw != "" {
		return cty.NilVal, false
	}
	return e.lit, true
}

// empty reports whether the expression carries nothing worth writing.
func (e *expr) empty() bool {
	switch {
	case e == nil:
		return true
	case e.isObj:
		return len(e.obj) == 0
	case e.isList:
		return len(e.list) == 0
	case e.raw != "":
		return false
	default:
		return e.lit == cty.NilVal || e.lit.IsNull()
	}
}

// tokens renders the expression as hclwrite tokens. Formatting is left to hclwrite.
func (e *expr) tokens() hclwrite.Tokens {
	switch {
	case e.isObj:
		tokens := hclwrite.Tokens{
			{Type: hclsyntax.TokenOBrace, Bytes: []byte("{")},
			{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		}
		for _, name := range sortedKeys(e.obj) {
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenIdent, Bytes: []byte(name)})
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenEqual, Bytes: []byte("=")})
			tokens = append(tokens, e.obj[name].tokens()...)
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")})
		}
		return append(tokens, &hclwrite.Token{Type: hclsyntax.TokenCBrace, Bytes: []byte("}")})
	case e.isList:
		tokens := hclwrite.Tokens{
			{Type: hclsyntax.TokenOBrack, Bytes: []byte("[")},
			{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		}
		for _, item := range e.list {
			tokens = append(tokens, item.tokens()...)
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenComma, Bytes: []byte(",")})
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")})
		}
		return append(tokens, &hclwrite.Token{Type: hclsyntax.TokenCBrack, Bytes: []byte("]")})
	case e.raw != "":
		return hclwrite.Tokens{{Type: hclsyntax.TokenIdent, Bytes: []byte(e.raw)}}
	default:
		return hclwrite.TokensForValue(e.lit)
	}
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// ctyTypeFromTF converts a Terraform protocol type into the cty type hclwrite renders values with. Both sides
// describe the same value space, so the conversion is mechanical and total for the types a schema can produce.
func ctyTypeFromTF(ty tftypes.Type) (cty.Type, error) {
	switch {
	case ty.Is(tftypes.String):
		return cty.String, nil
	case ty.Is(tftypes.Bool):
		return cty.Bool, nil
	case ty.Is(tftypes.Number):
		return cty.Number, nil
	case ty.Is(tftypes.DynamicPseudoType):
		return cty.DynamicPseudoType, nil
	}
	switch value := ty.(type) {
	case tftypes.List:
		element, err := ctyTypeFromTF(value.ElementType)
		if err != nil {
			return cty.NilType, err
		}
		return cty.List(element), nil
	case tftypes.Set:
		element, err := ctyTypeFromTF(value.ElementType)
		if err != nil {
			return cty.NilType, err
		}
		return cty.Set(element), nil
	case tftypes.Map:
		element, err := ctyTypeFromTF(value.ElementType)
		if err != nil {
			return cty.NilType, err
		}
		return cty.Map(element), nil
	case tftypes.Tuple:
		elements := make([]cty.Type, 0, len(value.ElementTypes))
		for _, elementType := range value.ElementTypes {
			element, err := ctyTypeFromTF(elementType)
			if err != nil {
				return cty.NilType, err
			}
			elements = append(elements, element)
		}
		return cty.Tuple(elements), nil
	case tftypes.Object:
		attributes := map[string]cty.Type{}
		for name, attributeType := range value.AttributeTypes {
			attribute, err := ctyTypeFromTF(attributeType)
			if err != nil {
				return cty.NilType, err
			}
			attributes[name] = attribute
		}
		return cty.Object(attributes), nil
	}
	return cty.NilType, fmt.Errorf("unsupported attribute type %s", ty.String())
}

// valueToCty converts a value decoded from Terraform state JSON into the given cty type, failing on any mismatch so
// that a legacy attribute whose shape changed is reported instead of silently reinterpreted.
func valueToCty(value any, ty cty.Type, path string) (cty.Value, error) {
	if value == nil {
		return cty.NullVal(ty), nil
	}
	switch {
	case ty == cty.String:
		text, ok := value.(string)
		if !ok {
			return cty.NilVal, fmt.Errorf("%s: expected a string, found %s", path, jsonKind(value))
		}
		return cty.StringVal(text), nil
	case ty == cty.Bool:
		flag, ok := value.(bool)
		if !ok {
			return cty.NilVal, fmt.Errorf("%s: expected a boolean, found %s", path, jsonKind(value))
		}
		return cty.BoolVal(flag), nil
	case ty == cty.Number:
		switch number := value.(type) {
		case json.Number:
			converted, err := cty.ParseNumberVal(number.String())
			if err != nil {
				return cty.NilVal, fmt.Errorf("%s: invalid number %q", path, number.String())
			}
			return converted, nil
		case float64:
			return cty.NumberVal(big.NewFloat(number)), nil
		default:
			return cty.NilVal, fmt.Errorf("%s: expected a number, found %s", path, jsonKind(value))
		}
	case ty.IsListType(), ty.IsSetType(), ty.IsTupleType():
		items, ok := value.([]any)
		if !ok {
			return cty.NilVal, fmt.Errorf("%s: expected a list, found %s", path, jsonKind(value))
		}
		return collectionToCty(items, ty, path)
	case ty.IsMapType():
		entries, ok := value.(map[string]any)
		if !ok {
			return cty.NilVal, fmt.Errorf("%s: expected an object, found %s", path, jsonKind(value))
		}
		if len(entries) == 0 {
			return cty.MapValEmpty(ty.ElementType()), nil
		}
		values := map[string]cty.Value{}
		for _, key := range sortedKeys(entries) {
			element, err := valueToCty(entries[key], ty.ElementType(), path+"."+key)
			if err != nil {
				return cty.NilVal, err
			}
			values[key] = element
		}
		return cty.MapVal(values), nil
	case ty.IsObjectType():
		entries, ok := value.(map[string]any)
		if !ok {
			return cty.NilVal, fmt.Errorf("%s: expected an object, found %s", path, jsonKind(value))
		}
		values := map[string]cty.Value{}
		for name, attributeType := range ty.AttributeTypes() {
			values[name] = cty.NullVal(attributeType)
			if entry, present := entries[name]; present {
				element, err := valueToCty(entry, attributeType, path+"."+name)
				if err != nil {
					return cty.NilVal, err
				}
				values[name] = element
			}
		}
		for _, key := range sortedKeys(entries) {
			if _, known := ty.AttributeTypes()[key]; !known {
				return cty.NilVal, fmt.Errorf("%s.%s: attribute has no destination", path, key)
			}
		}
		return cty.ObjectVal(values), nil
	}
	return cty.NilVal, fmt.Errorf("%s: unsupported destination type %s", path, ty.FriendlyName())
}

func collectionToCty(items []any, ty cty.Type, path string) (cty.Value, error) {
	if ty.IsTupleType() {
		types := ty.TupleElementTypes()
		if len(items) != len(types) {
			return cty.NilVal, fmt.Errorf("%s: expected %d elements, found %d", path, len(types), len(items))
		}
		values := make([]cty.Value, 0, len(items))
		for i, item := range items {
			element, err := valueToCty(item, types[i], fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return cty.NilVal, err
			}
			values = append(values, element)
		}
		return cty.TupleVal(values), nil
	}
	element := ty.ElementType()
	if len(items) == 0 {
		if ty.IsSetType() {
			return cty.SetValEmpty(element), nil
		}
		return cty.ListValEmpty(element), nil
	}
	values := make([]cty.Value, 0, len(items))
	for i, item := range items {
		converted, err := valueToCty(item, element, fmt.Sprintf("%s[%d]", path, i))
		if err != nil {
			return cty.NilVal, err
		}
		values = append(values, converted)
	}
	if ty.IsSetType() {
		return cty.SetVal(values), nil
	}
	return cty.ListVal(values), nil
}

func jsonKind(value any) string {
	switch value.(type) {
	case string:
		return "a string"
	case bool:
		return "a boolean"
	case float64, json.Number:
		return "a number"
	case []any:
		return "a list"
	case map[string]any:
		return "an object"
	default:
		return fmt.Sprintf("%T", value)
	}
}
