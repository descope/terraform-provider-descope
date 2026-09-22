package emit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type Resource struct {
	Type         string // the Terraform resource type, e.g. "descope_flow"
	Label        string // the unique resource label
	EntityID     string // the server-side entity ID, for cross-resource references
	ImportID     string // the ID for the resource's import block
	Scope        string // the method or app_id value for scoped resources
	HasProjectID bool
	HasMethod    bool
	HasAppID     bool
	Attrs        []prune.Attr
}

type Plan struct {
	ProjectID string
	Resources []Resource
}

func Write(ctx context.Context, outDir string, plan *Plan) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	slices.SortFunc(plan.Resources, func(a, b Resource) int {
		if c := strings.Compare(areaFile(a.Type), areaFile(b.Type)); c != 0 {
			return c
		}
		if c := strings.Compare(a.Type, b.Type); c != 0 {
			return c
		}
		return strings.Compare(a.Label, b.Label)
	})

	files := map[string]*hclwrite.File{}
	fileBody := func(name string) *hclwrite.Body {
		file, ok := files[name]
		if !ok {
			file = hclwrite.NewEmptyFile()
			files[name] = file
		} else {
			file.Body().AppendNewline()
		}
		return file.Body()
	}

	refs := references(plan)
	nameRefs := nameReferences(ctx, plan)
	projectRef := projectReference(plan)

	var variables []variable
	for _, resource := range plan.Resources {
		files, err := extractFiles(ctx, outDir, &resource)
		if err != nil {
			return fmt.Errorf("extracting files for %s.%s: %w", resource.Type, resource.Label, err)
		}
		body := fileBody(areaFile(resource.Type))
		block := body.AppendNewBlock("resource", []string{resource.Type, resource.Label})
		if err := resourceBody(ctx, block.Body(), resource, refs, nameRefs, projectRef, files, &variables); err != nil {
			return fmt.Errorf("emitting %s.%s: %w", resource.Type, resource.Label, err)
		}
	}

	writeProviderFile(fileBody("provider.tf"))
	writeVariablesFile(fileBody("variables.tf"), variables, projectRef == nil)
	writeImportsFile(fileBody("import.tf"), plan.Resources)

	for name, file := range files {
		content := hclwrite.Format(file.Bytes())
		if strings.TrimSpace(string(content)) == "" {
			continue // e.g. variables.tf when the project supplies the id and nothing sensitive was promoted
		}
		if err := os.WriteFile(filepath.Join(outDir, name), content, 0o644); err != nil {
			return err
		}
	}
	return nil
}

type variable struct {
	name        string
	description string
}

func resourceBody(ctx context.Context, body *hclwrite.Body, resource Resource, refs, nameRefs map[string]hcl.Traversal, projectRef hcl.Traversal, files map[string]string, variables *[]variable) error {
	if resource.HasProjectID {
		if projectRef == nil {
			projectRef = hcl.Traversal{hcl.TraverseRoot{Name: "var"}, hcl.TraverseAttr{Name: "project_id"}}
		}
		body.SetAttributeTraversal("project_id", projectRef)
	}
	if resource.HasMethod && resource.Scope != "" {
		body.SetAttributeValue("method", cty.StringVal(resource.Scope))
	}
	if resource.HasAppID && resource.Scope != "" {
		if traversal, ok := refs[resource.Scope]; ok {
			body.SetAttributeTraversal("app_id", traversal)
		} else {
			body.SetAttributeValue("app_id", cty.StringVal(resource.Scope))
		}
	}

	attrs := slices.Clone(resource.Attrs)
	slices.SortStableFunc(attrs, func(a, b prune.Attr) int {
		if (a.Name == "name") != (b.Name == "name") {
			if a.Name == "name" {
				return -1
			}
			return 1
		}
		return 0
	})

	for _, attr := range attrs {
		if relative, ok := files[attr.Name]; ok {
			body.SetAttributeRaw(attr.Name, hclwrite.Tokens{
				{Type: hclsyntax.TokenIdent, Bytes: []byte(fileTokens(relative))},
			})
			continue
		}
		if resource.Type == "descope_role" && attr.Name == "permissions" {
			if tokens, ok := permissionListTokens(ctx, attr, nameRefs); ok {
				body.SetAttributeRaw(attr.Name, tokens)
				continue
			}
		}
		tokens, err := attrTokens(ctx, attr, resource, refs, resource.Label, variables)
		if err != nil {
			return err
		}
		body.SetAttributeRaw(attr.Name, tokens)
	}
	return nil
}

func attrTokens(ctx context.Context, attr prune.Attr, resource Resource, refs map[string]hcl.Traversal, prefix string, variables *[]variable) (hclwrite.Tokens, error) {
	if attr.Placeholder {
		name := sanitizeLabel(prefix + "_" + attr.Name)
		*variables = append(*variables, variable{
			name:        name,
			description: fmt.Sprintf("Secret value for the %s attribute of %s.%s", attr.Name, resource.Type, resource.Label),
		})
		return hclwrite.TokensForTraversal(hcl.Traversal{hcl.TraverseRoot{Name: "var"}, hcl.TraverseAttr{Name: name}}), nil
	}
	if attr.IsNested {
		return objectTokens(ctx, attr.Nested, resource, refs, prefix+"_"+attr.Name, variables)
	}
	if attr.IsElements {
		elements := make([]hclwrite.Tokens, 0, len(attr.Elements))
		for i, element := range attr.Elements {
			// a single element needs no index, and most lists have one; more than one and the index is what keeps their variables apart
			elementPrefix := prefix + "_" + attr.Name
			if len(attr.Elements) > 1 {
				elementPrefix = fmt.Sprintf("%s_%d", elementPrefix, i)
			}
			tokens, err := objectTokens(ctx, element, resource, refs, elementPrefix, variables)
			if err != nil {
				return nil, err
			}
			elements = append(elements, tokens)
		}
		return hclwrite.TokensForTuple(elements), nil
	}
	if s, ok := stringValue(ctx, attr.Value); ok && s != resource.EntityID { // no self-references
		if traversal, ok := refs[s]; ok {
			return hclwrite.TokensForTraversal(traversal), nil
		}
	}
	value, err := ctyValue(ctx, attr.Value)
	if err != nil {
		return nil, err
	}
	return hclwrite.TokensForValue(value), nil
}

func objectTokens(ctx context.Context, attrs []prune.Attr, resource Resource, refs map[string]hcl.Traversal, prefix string, variables *[]variable) (hclwrite.Tokens, error) {
	items := make([]hclwrite.ObjectAttrTokens, 0, len(attrs))
	for _, attr := range attrs {
		tokens, err := attrTokens(ctx, attr, resource, refs, prefix, variables)
		if err != nil {
			return nil, err
		}
		items = append(items, hclwrite.ObjectAttrTokens{
			Name:  hclwrite.TokensForIdentifier(attr.Name),
			Value: tokens,
		})
	}
	return hclwrite.TokensForObject(items), nil
}

func writeProviderFile(body *hclwrite.Body) {
	terraform := body.AppendNewBlock("terraform", nil)
	providers := terraform.Body().AppendNewBlock("required_providers", nil)
	providers.Body().SetAttributeValue("descope", cty.ObjectVal(map[string]cty.Value{
		"source": cty.StringVal("descope/descope"),
	}))
	body.AppendNewline()
	provider := body.AppendNewBlock("provider", []string{"descope"})
	provider.Body().AppendUnstructuredTokens(hclwrite.Tokens{
		{Type: hclsyntax.TokenComment, Bytes: []byte("# authenticates with the DESCOPE_MANAGEMENT_KEY environment variable\n")},
	})
}

func writeVariablesFile(body *hclwrite.Body, variables []variable, declareProjectID bool) {
	first := true
	if declareProjectID {
		block := body.AppendNewBlock("variable", []string{"project_id"})
		block.Body().SetAttributeValue("description", cty.StringVal("The ID of the Descope project"))
		block.Body().SetAttributeTraversal("type", hcl.Traversal{hcl.TraverseRoot{Name: "string"}})
		first = false
	}

	slices.SortFunc(variables, func(a, b variable) int { return strings.Compare(a.name, b.name) })
	variables = slices.CompactFunc(variables, func(a, b variable) bool { return a.name == b.name })
	for _, v := range variables {
		if !first {
			body.AppendNewline()
		}
		first = false
		block := body.AppendNewBlock("variable", []string{v.name})
		block.Body().SetAttributeValue("description", cty.StringVal(v.description))
		block.Body().SetAttributeTraversal("type", hcl.Traversal{hcl.TraverseRoot{Name: "string"}})
		block.Body().SetAttributeValue("sensitive", cty.True)
	}
}

func writeImportsFile(body *hclwrite.Body, resources []Resource) {
	for i, resource := range resources {
		if i > 0 {
			body.AppendNewline()
		}
		block := body.AppendNewBlock("import", nil)
		block.Body().SetAttributeTraversal("to", hcl.Traversal{
			hcl.TraverseRoot{Name: resource.Type},
			hcl.TraverseAttr{Name: resource.Label},
		})
		block.Body().SetAttributeValue("id", cty.StringVal(resource.ImportID))
	}
}

func areaFile(resourceType string) string {
	if strings.HasSuffix(resourceType, "_connector") {
		return "connectors.tf"
	}
	switch resourceType {
	case "descope_project":
		return "project.tf"
	case "descope_flow", "descope_widget", "descope_styles":
		return "flows.tf"
	case "descope_email_template", "descope_text_template", "descope_jwt_template":
		return "templates.tf"
	case "descope_oidc_app", "descope_saml_app", "descope_wsfed_app", "descope_inbound_app", "descope_outbound_app", "descope_app_role", "descope_app_permission":
		return "apps.tf"
	case "descope_role", "descope_permission":
		return "authorization.tf"
	case "descope_admin_portal", "descope_fga_schema", "descope_session_migration":
		return "auth.tf"
	}
	if strings.HasSuffix(resourceType, "_settings") {
		return "auth.tf"
	}
	return "misc.tf"
}

func permissionListTokens(ctx context.Context, attr prune.Attr, nameRefs map[string]hcl.Traversal) (hclwrite.Tokens, bool) {
	value, err := ctyValue(ctx, attr.Value)
	if err != nil || value.IsNull() || !value.CanIterateElements() {
		return nil, false
	}
	var elements []hclwrite.Tokens
	substituted := false
	for iterator := value.ElementIterator(); iterator.Next(); {
		_, element := iterator.Element()
		if element.Type() != cty.String || element.IsNull() {
			return nil, false
		}
		if traversal, ok := nameRefs[element.AsString()]; ok {
			elements = append(elements, hclwrite.TokensForTraversal(traversal))
			substituted = true
		} else {
			elements = append(elements, hclwrite.TokensForValue(element))
		}
	}
	if !substituted {
		return nil, false
	}
	return hclwrite.TokensForTuple(elements), true
}
