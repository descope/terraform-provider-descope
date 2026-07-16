package outboundapp

import (
	"context"
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func Test_Model_round_trips_outbound_app_configuration(t *testing.T) {
	// Given
	model := OutboundAppModel{
		ID:                     types.StringValue("app-1"),
		Name:                   types.StringValue("example"),
		Description:            types.StringValue("description"),
		ClientID:               types.StringValue("client-1"),
		Logo:                   types.StringValue("https://example.com/logo.png"),
		DiscoveryURL:           types.StringValue("https://example.com/.well-known/openid-configuration"),
		AuthorizationURL:       types.StringValue("https://example.com/oauth/authorize"),
		AuthorizationURLParams: urlParamsValue(context.Background(), []descope.URLParam{{Key: "audience", Value: "api"}}),
		TokenURL:               types.StringValue("https://example.com/oauth/token"),
		TokenURLParams:         urlParamsValue(context.Background(), []descope.URLParam{{Key: "resource", Value: "api"}}),
		RevocationURL:          types.StringValue("https://example.com/oauth/revoke"),
		DefaultScopes:          stringSetValue(context.Background(), []string{"openid", "profile"}),
		DefaultRedirectURL:     types.StringValue("https://example.com/callback"),
		CallbackDomain:         types.StringValue("https://example.com"),
		PKCE:                   types.BoolValue(true),
		AccessType:             types.StringValue("offline"),
		Prompt:                 stringSetValue(context.Background(), []string{"consent"}),
	}
	var diagnostics diag.Diagnostics

	// When
	app := model.OutboundApp(context.Background(), &diagnostics)
	var roundTrip OutboundAppModel
	roundTrip.SetOutboundApp(context.Background(), app)

	// Then
	require.False(t, diagnostics.HasError())
	require.Equal(t, model.Name, roundTrip.Name)
	require.Equal(t, model.AuthorizationURLParams, roundTrip.AuthorizationURLParams)
	require.Equal(t, model.DefaultScopes, roundTrip.DefaultScopes)
	require.Equal(t, model.Prompt, roundTrip.Prompt)
}

func Test_Schema_marks_client_secret_write_only(t *testing.T) {
	// When
	attribute, ok := Schema.Attributes["client_secret"].(schema.StringAttribute)

	// Then
	require.True(t, ok)
	require.True(t, attribute.Optional)
	require.True(t, attribute.Sensitive)
	require.True(t, attribute.WriteOnly)
	require.False(t, attribute.Computed)
}
