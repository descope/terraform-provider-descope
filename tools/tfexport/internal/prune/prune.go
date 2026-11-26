package prune

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Attr struct {
	Name        string
	Value       attr.Value // a leaf value or a whole plain collection
	Nested      []Attr     // the surviving fields of a single-nested object
	Elements    [][]Attr   // the elements of a nested list or set, each with its surviving fields
	IsNested    bool       // set when Nested is the meaningful field, since it may be empty
	IsElements  bool       // set when Elements is the meaningful field, since it may be empty
	Placeholder bool       // a required sensitive attribute: emit var.<label>_<name>
	Verbatim    bool       // an all-default object emitted only to keep the plan stable
}

func Meaningful(attrs []Attr) bool {
	for _, a := range attrs {
		if !a.Verbatim {
			return true
		}
	}
	return false
}

type Secret struct {
	Path     string // the attribute path within the resource, e.g. "api_key"
	Required bool   // the attribute must be assigned a variable for the config to apply
	Dropped  bool   // the backend returned a value that the configuration cannot carry, so applying it clears the secret
}

// scoping attributes are re-added by the emitter. id is absent on purpose: generated ids drop themselves, user-chosen ids are required.
var rootAttributes = []string{"project_id", "method", "app_id"}

func Object(ctx context.Context, attributes map[string]schema.Attribute, object types.Object, root bool) (attrs []Attr, secrets []Secret) {
	values := object.Attributes()

	names := make([]string, 0, len(attributes))
	for name := range attributes {
		names = append(names, name)
	}
	slices.Sort(names)

	for _, name := range names {
		attribute := attributes[name]
		value, ok := values[name]
		if !ok {
			continue
		}
		if root && slices.Contains(rootAttributes, name) {
			continue
		}
		if attribute.IsComputed() && !attribute.IsOptional() && !attribute.IsRequired() {
			continue // computed-only attributes cannot be set in configurations
		}

		if attribute.IsSensitive() {
			switch {
			case attribute.IsRequired():
				attrs = append(attrs, Attr{Name: name, Placeholder: true})
				secrets = append(secrets, Secret{Path: name, Required: true})
			case isMap(attribute) && !value.IsNull() && !isServerZeroValue(ctx, attribute, value):
				nested := secretMapKeys(ctx, value, name, &secrets)
				attrs = append(attrs, Attr{Name: name, Nested: nested, IsNested: true})
			case !value.IsNull() && !isServerZeroValue(ctx, attribute, value):
				if _, hasDefault := AttributeDefault(ctx, attribute); hasDefault {
					attrs = append(attrs, Attr{Name: name, Placeholder: true})
					secrets = append(secrets, Secret{Path: name, Required: true})
				} else if attribute.IsComputed() {
					secrets = append(secrets, Secret{Path: name})
				} else {
					secrets = append(secrets, Secret{Path: name, Dropped: true})
				}
			default:
				secrets = append(secrets, Secret{Path: name})
			}
			continue
		}

		if value.IsNull() && !attribute.IsRequired() {
			continue // never set, and matches the imported state
		}

		if def, ok := AttributeDefault(ctx, attribute); ok {
			if valuesEqual(ctx, value, def) {
				continue
			}
		} else if isServerZeroValue(ctx, attribute, value) {
			continue // the backend fills these in rather than leaving them null, and Computed keeps the stored value
		}

		switch a := attribute.(type) {
		case schema.SingleNestedAttribute:
			nested := nestedObject(ctx, a.Attributes, value, name, &secrets)
			verbatim := false
			if len(nested) == 0 {
				// every field matches its nested default but the object differs from its own, so the non-null fields are emitted verbatim
				nested = verbatimObject(ctx, a.Attributes, value)
				verbatim = true
			}
			if len(nested) == 0 && !attribute.IsRequired() {
				continue
			}
			attrs = append(attrs, Attr{Name: name, Nested: nested, IsNested: true, Verbatim: verbatim})
		case schema.ListNestedAttribute:
			attrs = append(attrs, elementsAttr(ctx, a.NestedObject.Attributes, value, name, &secrets))
		case schema.SetNestedAttribute:
			attrs = append(attrs, elementsAttr(ctx, a.NestedObject.Attributes, value, name, &secrets))
		default:
			attrs = append(attrs, Attr{Name: name, Value: value})
		}
	}

	return attrs, secrets
}

func nestedObject(ctx context.Context, attributes map[string]schema.Attribute, value attr.Value, path string, secrets *[]Secret) []Attr {
	valuable, ok := value.(basetypes.ObjectValuable)
	if !ok {
		return nil
	}
	object, diags := valuable.ToObjectValue(ctx)
	if diags.HasError() {
		return nil
	}
	nested, nestedSecrets := Object(ctx, attributes, object, false)
	for _, secret := range nestedSecrets {
		*secrets = append(*secrets, Secret{Path: path + "." + secret.Path, Required: secret.Required, Dropped: secret.Dropped})
	}
	return nested
}

func isMap(attribute schema.Attribute) bool {
	_, ok := attribute.(schema.MapAttribute)
	return ok
}

func secretMapKeys(ctx context.Context, value attr.Value, path string, secrets *[]Secret) []Attr {
	valuable, ok := value.(basetypes.MapValuable)
	if !ok {
		return nil
	}
	m, diags := valuable.ToMapValue(ctx)
	if diags.HasError() {
		return nil
	}
	keys := make([]string, 0, len(m.Elements()))
	for key := range m.Elements() {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	attrs := make([]Attr, 0, len(keys))
	for _, key := range keys {
		attrs = append(attrs, Attr{Name: key, Placeholder: true})
		*secrets = append(*secrets, Secret{Path: fmt.Sprintf("%s[%q]", path, key), Required: true})
	}
	return attrs
}

func verbatimObject(ctx context.Context, attributes map[string]schema.Attribute, value attr.Value) []Attr {
	valuable, ok := value.(basetypes.ObjectValuable)
	if !ok {
		return nil
	}
	object, diags := valuable.ToObjectValue(ctx)
	if diags.HasError() {
		return nil
	}
	values := object.Attributes()

	names := make([]string, 0, len(attributes))
	for name := range attributes {
		names = append(names, name)
	}
	slices.Sort(names)

	var attrs []Attr
	for _, name := range names {
		attribute := attributes[name]
		value, ok := values[name]
		if !ok || value.IsNull() || attribute.IsSensitive() {
			continue
		}
		if attribute.IsComputed() && !attribute.IsOptional() && !attribute.IsRequired() {
			continue
		}
		attrs = append(attrs, Attr{Name: name, Value: value})
	}
	return attrs
}

func elementsAttr(ctx context.Context, attributes map[string]schema.Attribute, value attr.Value, path string, secrets *[]Secret) Attr {
	result := Attr{Name: path, IsElements: true}

	var elements []attr.Value
	switch valuable := value.(type) {
	case basetypes.ListValuable:
		if list, diags := valuable.ToListValue(ctx); !diags.HasError() {
			elements = list.Elements()
		}
	case basetypes.SetValuable:
		if set, diags := valuable.ToSetValue(ctx); !diags.HasError() {
			elements = set.Elements()
		}
	}

	for i, element := range elements {
		result.Elements = append(result.Elements, nestedObject(ctx, attributes, element, fmt.Sprintf("%s[%d]", path, i), secrets))
	}
	return result
}
