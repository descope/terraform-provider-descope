package provider_test

import (
	"context"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Attributes allowed to carry RequiresReplace even though isPlannedReplace cannot detect them, keyed as "attribute" (any resource) or
// "resource_type attribute". DANGER: changing an allowlisted attribute recreates the resource, and if that resource has (or gains)
// deletion_protection the replacement will NOT be blocked. Only allowlist immutable identity attributes of unprotected resources.
var replaceModifierAllowlist = map[string]bool{
	"project_id":                         true, // moving a resource between projects is a recreate by definition
	"descope_app_permission app_id":      true,
	"descope_app_role app_id":            true,
	"descope_access_key bound_user_id":   true,
	"descope_access_key expire_time":     true,
	"descope_management_key expire_time": true,
	"descope_management_key rebac":       true,
	"descope_flow flow_id":               true,
	"descope_widget widget_id":           true,
	"descope_oauth_provider id":          true,
	"descope_email_template method":      true,
	"descope_text_template method":       true,
	"descope_voice_template method":      true,
	"descope_user_attribute id":          true,
	"descope_user_attribute type":        true,
	"descope_tenant_attribute id":        true,
	"descope_tenant_attribute type":      true,
	"descope_access_key_attribute id":    true,
	"descope_access_key_attribute type":  true,
}

// Guards the deletion protection invariant: isPlannedReplace only detects RequiresReplace on top level string and bool attributes, so
// unless allowlisted above no other attribute may carry the modifier in a protected resource, and none at all in an unprotected one.
func TestRequiresReplaceModifiers(t *testing.T) {
	ctx := context.Background()
	found := 0
	for _, newResource := range provider.NewDescopeProvider("test")().Resources(ctx) {
		r := newResource()
		metadata := resource.MetadataResponse{}
		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "descope"}, &metadata)
		response := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &response)
		require.Empty(t, response.Schema.Blocks, "update this test to check for RequiresReplace modifiers in the blocks of %s", metadata.TypeName)
		_, protected := response.Schema.Attributes["deletion_protection"]
		for name, attribute := range response.Schema.Attributes {
			found += checkReplaceModifiers(t, metadata.TypeName, name, attribute, protected, true)
		}
	}
	require.NotZero(t, found, "expected at least one RequiresReplace modifier to be detected, the detection heuristic is probably broken")
}

func checkReplaceModifiers(t *testing.T, resourceName, attributePath string, attribute schema.Attribute, protected, topLevel bool) int {
	found := 0
	if hasReplaceModifier(attribute) {
		found++
		_, isString := attribute.(schema.StringAttribute)
		_, isBool := attribute.(schema.BoolAttribute)
		detected := protected && topLevel && (isString || isBool)
		allowed := topLevel && (replaceModifierAllowlist[attributePath] || replaceModifierAllowlist[resourceName+" "+attributePath])
		if !detected && !allowed {
			assert.Fail(t, "Undetectable modifier", "the RequiresReplace modifier on %s in %s will not be detected by isPlannedReplace", attributePath, resourceName)
		}
	}
	nested := map[string]schema.Attribute{}
	switch a := attribute.(type) {
	case schema.SingleNestedAttribute:
		nested = a.Attributes
	case schema.ListNestedAttribute:
		nested = a.NestedObject.Attributes
	case schema.SetNestedAttribute:
		nested = a.NestedObject.Attributes
	case schema.MapNestedAttribute:
		nested = a.NestedObject.Attributes
	}
	for name, child := range nested {
		found += checkReplaceModifiers(t, resourceName, attributePath+"."+name, child, protected, false)
	}
	return found
}

func hasReplaceModifier(attribute schema.Attribute) bool {
	modifiers := reflect.ValueOf(attribute).FieldByName("PlanModifiers")
	for i := 0; modifiers.IsValid() && i < modifiers.Len(); i++ {
		if strings.Contains(reflect.TypeOf(modifiers.Index(i).Interface()).String(), "requiresReplace") {
			return true
		}
	}
	return false
}

// Every attribute holding a JSON document, as "resource_type attribute". Without the normalization modifier these plan a change
// whenever the configured document is formatted differently to the one in state, which is what every tfexport run produces.
var jsonAttributes = []string{
	"descope_access_key custom_attributes",
	"descope_access_key custom_claims",
	"descope_flow data",
	"descope_jwt_template template",
	"descope_list json",
	"descope_styles data",
	"descope_widget data",
}

// Guards that JSON attributes compare by content, not formatting: the listed ones must carry the normalization modifier and any
// attribute validated as JSON must be listed, so a new JSON attribute declared with stringattr instead of jsonattr fails here.
func TestJSONAttributeNormalization(t *testing.T) {
	ctx := context.Background()
	normalized := []string{}
	for _, newResource := range provider.NewDescopeProvider("test")().Resources(ctx) {
		r := newResource()
		metadata := resource.MetadataResponse{}
		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "descope"}, &metadata)
		response := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &response)
		for name, attribute := range response.Schema.Attributes {
			path := metadata.TypeName + " " + name
			if hasNamedModifier(attribute, "useStateWhenEquivalent") {
				normalized = append(normalized, path)
			} else if hasNamedValidator(attribute, "jsonValidator") {
				assert.Fail(t, "Unnormalized JSON attribute", "%s is validated as JSON but will plan a change on reformatting, declare it with jsonattr", path)
			}
		}
	}
	slices.Sort(normalized)
	assert.Equal(t, jsonAttributes, normalized)
}

func hasNamedModifier(attribute schema.Attribute, name string) bool {
	return hasNamedElement(attribute, "PlanModifiers", name)
}

func hasNamedValidator(attribute schema.Attribute, name string) bool {
	return hasNamedElement(attribute, "Validators", name)
}

func hasNamedElement(attribute schema.Attribute, field, name string) bool {
	elements := reflect.ValueOf(attribute).FieldByName(field)
	for i := 0; elements.IsValid() && i < elements.Len(); i++ {
		if strings.Contains(reflect.TypeOf(elements.Index(i).Interface()).String(), name) {
			return true
		}
	}
	return false
}
