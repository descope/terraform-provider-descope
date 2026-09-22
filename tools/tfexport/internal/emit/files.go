package emit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const inlineBodyLimit = 1024

func extractFiles(ctx context.Context, outDir string, resource *Resource) (map[string]string, error) {
	extracted := map[string]string{}

	extract := func(name, dir, filename string, pretty bool) error {
		attribute := findAttr(resource.Attrs, name)
		if attribute == nil {
			return nil
		}
		value, ok := stringValue(ctx, attribute.Value)
		if !ok || value == "" {
			return nil
		}
		content := []byte(value)
		if pretty {
			indented := &bytes.Buffer{}
			if err := json.Indent(indented, content, "", "  "); err == nil {
				indented.WriteByte('\n')
				content = indented.Bytes()
			}
		}
		if err := os.MkdirAll(filepath.Join(outDir, dir), 0o755); err != nil {
			return err
		}
		relative := filepath.Join(dir, filename)
		if err := os.WriteFile(filepath.Join(outDir, relative), content, 0o644); err != nil {
			return err
		}
		extracted[name] = filepath.ToSlash(relative)
		return nil
	}

	switch resource.Type {
	case "descope_flow":
		return extracted, extract("data", "flows", resource.Label+".json", true)
	case "descope_widget":
		return extracted, extract("data", "widgets", resource.Label+".json", true)
	case "descope_styles":
		return extracted, extract("data", ".", "styles.json", true)
	case "descope_email_template", "descope_text_template":
		for name, extension := range map[string]string{"html_body": ".html", "plain_text_body": ".txt", "body": ".txt"} {
			attribute := findAttr(resource.Attrs, name)
			if attribute == nil {
				continue
			}
			value, ok := stringValue(ctx, attribute.Value)
			if !ok || (len(value) <= inlineBodyLimit && !strings.Contains(value, "\n")) {
				continue
			}
			filename := resource.Label + extension
			if name == "plain_text_body" {
				filename = resource.Label + "_plain" + extension
			}
			if err := extract(name, "templates", filename, false); err != nil {
				return extracted, err
			}
		}
	}
	return extracted, nil
}

func findAttr(attrs []prune.Attr, name string) *prune.Attr {
	for i := range attrs {
		if attrs[i].Name == name && !attrs[i].IsNested && !attrs[i].IsElements && !attrs[i].Placeholder {
			return &attrs[i]
		}
	}
	return nil
}

func stringValue(ctx context.Context, value attr.Value) (string, bool) {
	valuable, ok := value.(basetypes.StringValuable)
	if !ok {
		return "", false
	}
	converted, diags := valuable.ToStringValue(ctx)
	if diags.HasError() || converted.IsNull() {
		return "", false
	}
	return converted.ValueString(), true
}

func fileTokens(relative string) string {
	return fmt.Sprintf(`file("${path.module}/%s")`, relative)
}
