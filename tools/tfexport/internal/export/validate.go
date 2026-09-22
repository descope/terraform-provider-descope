package export

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/resources"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func ensureValidConfig(ctx context.Context, exportable resources.ExportableResource, name string, original types.Object, attrs []prune.Attr, secrets []prune.Secret, warn func(format string, args ...any)) ([]prune.Attr, []prune.Secret) {
	schema := exportable.ExportSchema()
	value, err := configValue(ctx, schema.Attributes, original, attrs)
	if err != nil {
		warn("Failed to check %s %s: %s", exportable.ExportName(), name, err.Error())
		return attrs, secrets
	}
	initial := exportable.ExportValidate(ctx, value)
	if !initial.HasError() {
		return attrs, secrets
	}

	var candidates []string
	for _, secret := range secrets {
		if secret.Required {
			continue
		}
		candidates = append(candidates, secret.Path)
	}
	if len(candidates) == 0 {
		warn("The exported %s %s needs manual fixes: %s", exportable.ExportName(), name, validationDetails(initial))
		return attrs, secrets
	}

	validate := func(promoted []string) diag.Diagnostics {
		hypothetical := attrs
		for _, path := range promoted {
			hypothetical = withPlaceholder(hypothetical, path)
		}
		value, err := configValue(ctx, schema.Attributes, original, hypothetical)
		if err != nil {
			return diag.Diagnostics{diag.NewErrorDiagnostic("configuration", err.Error())}
		}
		return exportable.ExportValidate(ctx, value)
	}
	valid := func(promoted []string) bool { return !validate(promoted).HasError() }

	// what promoting every candidate secret still can't fix is what a person has to supply by hand
	if remaining := validate(candidates); remaining.HasError() {
		warn("The exported %s %s needs manual fixes: %s", exportable.ExportName(), name, validationDetails(remaining))
		return attrs, secrets
	}

	needed := slices.Clone(candidates)
	for _, name := range candidates {
		without := slices.DeleteFunc(slices.Clone(needed), func(s string) bool { return s == name })
		if valid(without) {
			needed = without
		}
	}

	for _, path := range needed {
		attrs = withPlaceholder(attrs, path)
	}
	updated := make([]prune.Secret, 0, len(secrets))
	for _, secret := range secrets {
		if slices.Contains(needed, secret.Path) {
			secret.Required = true
		}
		updated = append(updated, secret)
	}
	return attrs, updated
}

func validationDetails(diags diag.Diagnostics) string {
	details := make([]string, 0, len(diags.Errors()))
	for _, d := range diags.Errors() {
		details = append(details, d.Detail())
	}
	return strings.Join(details, "; ")
}

func withPlaceholder(attrs []prune.Attr, path string) []prune.Attr {
	segment, rest, nested := strings.Cut(path, ".")
	name, index, indexed := cutIndex(segment)
	result := slices.Clone(attrs)
	for i, attr := range result {
		if attr.Name != name {
			continue
		}
		switch {
		case indexed:
			if index >= len(attr.Elements) {
				return result // the tree doesn't describe this element, so there is nothing to promote
			}
			elements := slices.Clone(attr.Elements)
			elements[index] = withPlaceholder(elements[index], rest)
			result[i].IsElements = true
			result[i].Elements = elements
		case nested:
			result[i].IsNested = true
			result[i].Nested = withPlaceholder(attr.Nested, rest)
		default:
			result[i].Placeholder = true
		}
		return result
	}
	if indexed {
		return result // pruning always emits the collection itself, so a missing one can't be reconstructed
	}
	if nested {
		return append(result, prune.Attr{Name: name, IsNested: true, Nested: withPlaceholder(nil, rest)})
	}
	return append(result, prune.Attr{Name: name, Placeholder: true})
}

func cutIndex(segment string) (string, int, bool) {
	open := strings.IndexByte(segment, '[')
	if open < 0 || !strings.HasSuffix(segment, "]") {
		return segment, 0, false
	}
	index, err := strconv.Atoi(segment[open+1 : len(segment)-1])
	if err != nil || index < 0 {
		return segment, 0, false
	}
	return segment[:open], index, true
}

func configValue(ctx context.Context, attributes map[string]schema.Attribute, original types.Object, attrs []prune.Attr) (types.Object, error) {
	raw, err := original.ToTerraformValue(ctx)
	if err != nil {
		return types.Object{}, err
	}
	rebuilt, err := objectValue(ctx, attributes, raw, attrs)
	if err != nil {
		return types.Object{}, err
	}
	converted, err := original.Type(ctx).ValueFromTerraform(ctx, rebuilt)
	if err != nil {
		return types.Object{}, err
	}
	object, ok := converted.(types.Object)
	if !ok {
		return types.Object{}, fmt.Errorf("unexpected converted type %T", converted)
	}
	return object, nil
}

func objectValue(ctx context.Context, attributes map[string]schema.Attribute, raw tftypes.Value, attrs []prune.Attr) (tftypes.Value, error) {
	objectType, ok := raw.Type().(tftypes.Object)
	if !ok {
		return tftypes.Value{}, fmt.Errorf("unexpected value type %s", raw.Type())
	}
	values := map[string]tftypes.Value{}
	if err := raw.As(&values); err != nil {
		return tftypes.Value{}, err
	}

	surviving := map[string]prune.Attr{}
	for _, attr := range attrs {
		surviving[attr.Name] = attr
	}

	modified := map[string]tftypes.Value{}
	for name, attributeType := range objectType.AttributeTypes {
		attr, ok := surviving[name]
		switch {
		case !ok:
			modified[name] = tftypes.NewValue(attributeType, nil)
			if attribute, ok := attributes[name]; ok {
				if def, ok := prune.AttributeDefault(ctx, attribute); ok {
					if converted, err := def.ToTerraformValue(ctx); err == nil {
						modified[name] = converted
					}
				}
			}
		case attr.Placeholder:
			modified[name] = tftypes.NewValue(attributeType, tftypes.UnknownValue)
		case attr.IsNested, attr.IsElements:
			nested, err := nestedValue(ctx, attributes[name], values[name], attr)
			if err != nil {
				return tftypes.Value{}, err
			}
			modified[name] = nested
		default:
			modified[name] = values[name]
		}
	}
	return tftypes.NewValue(objectType, modified), nil
}

func nestedValue(ctx context.Context, attribute schema.Attribute, raw tftypes.Value, attr prune.Attr) (tftypes.Value, error) {
	if raw.IsNull() || !raw.IsKnown() {
		return raw, nil
	}
	switch a := attribute.(type) {
	case schema.SingleNestedAttribute:
		return objectValue(ctx, a.Attributes, raw, attr.Nested)
	case schema.ListNestedAttribute:
		return elementValues(ctx, a.NestedObject.Attributes, raw, attr.Elements)
	case schema.SetNestedAttribute:
		return elementValues(ctx, a.NestedObject.Attributes, raw, attr.Elements)
	}
	return raw, nil
}

func elementValues(ctx context.Context, attributes map[string]schema.Attribute, raw tftypes.Value, elements [][]prune.Attr) (tftypes.Value, error) {
	var values []tftypes.Value
	if err := raw.As(&values); err != nil {
		return tftypes.Value{}, err
	}
	if len(values) != len(elements) {
		return raw, nil // pruning never drops elements, so a mismatch means the tree doesn't describe this value
	}
	rebuilt := make([]tftypes.Value, 0, len(values))
	for i, value := range values {
		element, err := objectValue(ctx, attributes, value, elements[i])
		if err != nil {
			return tftypes.Value{}, err
		}
		rebuilt = append(rebuilt, element)
	}
	return tftypes.NewValue(raw.Type(), rebuilt), nil
}
