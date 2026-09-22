package export

import (
	"context"
	"fmt"
	"io"
	"maps"
	"slices"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/read"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const redacted = "(redacted)"

// PrintState renders what was read, masking sensitive values: the API returns them in the clear and sensitivity is
// schema metadata rather than a property of the value, so nothing else would keep them out of this output.
func PrintState(ctx context.Context, w io.Writer, results []read.Result) {
	loaded := registry.Load(ctx)
	for _, result := range results {
		fmt.Fprintf(w, "### %s %q (id %s)\n", result.Instance.Resource, result.Instance.Name, result.Instance.ID) // nolint:errcheck
		var attributes map[string]schema.Attribute
		if exportable, ok := loaded[result.Instance.Resource]; ok {
			attributes = exportable.ExportSchema().Attributes
		}
		printObject(ctx, w, attributes, result.Object, "  ")
		fmt.Fprintln(w) // nolint:errcheck
	}
}

func printObject(ctx context.Context, w io.Writer, attributes map[string]schema.Attribute, object types.Object, indent string) {
	values := object.Attributes()
	for _, name := range slices.Sorted(maps.Keys(values)) {
		value := values[name]
		attribute := attributes[name]
		if attribute != nil && attribute.IsSensitive() && !value.IsNull() {
			fmt.Fprintf(w, "%s%s = %s\n", indent, name, redacted) // nolint:errcheck
			continue
		}

		var nested map[string]schema.Attribute
		switch a := attribute.(type) {
		case schema.SingleNestedAttribute:
			nested = a.Attributes
		case schema.ListNestedAttribute:
			nested = a.NestedObject.Attributes
		case schema.SetNestedAttribute:
			nested = a.NestedObject.Attributes
		}
		if nested == nil || value.IsNull() || value.IsUnknown() {
			fmt.Fprintf(w, "%s%s = %s\n", indent, name, value.String()) // nolint:errcheck
			continue
		}

		if _, ok := attribute.(schema.SingleNestedAttribute); ok {
			fmt.Fprintf(w, "%s%s = {\n", indent, name) // nolint:errcheck
			printNested(ctx, w, nested, value, indent+"  ")
			fmt.Fprintf(w, "%s}\n", indent) // nolint:errcheck
			continue
		}

		elements, ok := value.(interface{ Elements() []attr.Value })
		if !ok {
			fmt.Fprintf(w, "%s%s = %s\n", indent, name, value.String()) // nolint:errcheck
			continue
		}
		fmt.Fprintf(w, "%s%s = [\n", indent, name) // nolint:errcheck
		for _, element := range elements.Elements() {
			fmt.Fprintf(w, "%s  {\n", indent) // nolint:errcheck
			printNested(ctx, w, nested, element, indent+"    ")
			fmt.Fprintf(w, "%s  }\n", indent) // nolint:errcheck
		}
		fmt.Fprintf(w, "%s]\n", indent) // nolint:errcheck
	}
}

func printNested(ctx context.Context, w io.Writer, attributes map[string]schema.Attribute, value attr.Value, indent string) {
	valuable, ok := value.(basetypes.ObjectValuable)
	if !ok {
		fmt.Fprintf(w, "%s%s\n", indent, value.String()) // nolint:errcheck
		return
	}
	object, diags := valuable.ToObjectValue(ctx)
	if diags.HasError() {
		fmt.Fprintf(w, "%s%s\n", indent, value.String()) // nolint:errcheck
		return
	}
	printObject(ctx, w, attributes, object, indent)
}
