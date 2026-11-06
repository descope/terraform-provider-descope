package resources_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/descope/terraform-provider-descope/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

// connectorSchema returns the schema of the given provider resource type name.
func connectorSchema(t *testing.T, typeName string) schema.Schema {
	t.Helper()
	for _, constructor := range resources.ConnectorResources() {
		res := constructor()
		var md resource.MetadataResponse
		res.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "descope"}, &md)
		if md.TypeName == typeName {
			var sr resource.SchemaResponse
			res.Schema(context.Background(), resource.SchemaRequest{}, &sr)
			return sr.Schema
		}
	}
	t.Fatalf("resource %s not found", typeName)
	return schema.Schema{}
}

// Test that a read after an import populates schema defaults for attributes the server response omits.
func TestReadPopulatesSchemaDefaultsOmittedByServer(t *testing.T) {
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"C1","type":"external-token-http","name":"test","description":"","disabled":false,"configuration":{"endpoint":"https://example.com"}}`))
	}))
	defer server.Close()

	p, err := providerserver.NewProtocol6WithError(provider.NewDescopeProvider("test")())()
	require.NoError(t, err)

	configType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"project_id":     tftypes.String,
		"management_key": tftypes.String,
		"base_url":       tftypes.String,
	}}
	config, err := tfprotov6.NewDynamicValue(configType, tftypes.NewValue(configType, map[string]tftypes.Value{
		"project_id":     tftypes.NewValue(tftypes.String, nil),
		"management_key": tftypes.NewValue(tftypes.String, "test-key"),
		"base_url":       tftypes.NewValue(tftypes.String, server.URL),
	}))
	require.NoError(t, err)
	_, err = p.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{TerraformVersion: "1.13.0", Config: &config})
	require.NoError(t, err)

	sch := connectorSchema(t, "descope_external_token_http_connector")
	typ := sch.Type().TerraformType(ctx)
	attributes := map[string]tftypes.Value{}
	for name, attributeType := range typ.(tftypes.Object).AttributeTypes {
		attributes[name] = tftypes.NewValue(attributeType, nil)
	}
	attributes["id"] = tftypes.NewValue(tftypes.String, "C1")
	attributes["project_id"] = tftypes.NewValue(tftypes.String, "P1")
	state, err := tfprotov6.NewDynamicValue(typ, tftypes.NewValue(typ, attributes))
	require.NoError(t, err)

	read, err := p.ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName:     "descope_external_token_http_connector",
		CurrentState: &state,
	})
	require.NoError(t, err)
	for _, d := range read.Diagnostics {
		require.NotEqual(t, tfprotov6.DiagnosticSeverityError, d.Severity, d.Detail)
	}

	newState, err := read.NewState.Unmarshal(typ)
	require.NoError(t, err)
	var final map[string]tftypes.Value
	require.NoError(t, newState.As(&final))
	useStaticIPs := final["use_static_ips"]
	require.False(t, useStaticIPs.IsNull(), "use_static_ips should hold the schema default, not null")
	var value bool
	require.NoError(t, useStaticIPs.As(&value))
	require.False(t, value)
}
