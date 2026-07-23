package outboundapp

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/listtype"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/valuesettype"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type URLParamModel struct {
	Key   stringattr.Type `tfsdk:"key"`
	Value stringattr.Type `tfsdk:"value"`
}

type OutboundAppModel struct {
	ID                     stringattr.Type              `tfsdk:"id"`
	ProjectID              stringattr.Type              `tfsdk:"project_id"`
	Name                   stringattr.Type              `tfsdk:"name"`
	Description            stringattr.Type              `tfsdk:"description"`
	ClientID               stringattr.Type              `tfsdk:"client_id"`
	ClientSecret           stringattr.Type              `tfsdk:"client_secret"`
	Logo                   stringattr.Type              `tfsdk:"logo"`
	DiscoveryURL           stringattr.Type              `tfsdk:"discovery_url"`
	AuthorizationURL       stringattr.Type              `tfsdk:"authorization_url"`
	AuthorizationURLParams listattr.Type[URLParamModel] `tfsdk:"authorization_url_params"`
	TokenURL               stringattr.Type              `tfsdk:"token_url"`
	TokenURLParams         listattr.Type[URLParamModel] `tfsdk:"token_url_params"`
	RevocationURL          stringattr.Type              `tfsdk:"revocation_url"`
	DefaultScopes          strsetattr.Type              `tfsdk:"default_scopes"`
	DefaultRedirectURL     stringattr.Type              `tfsdk:"default_redirect_url"`
	CallbackDomain         stringattr.Type              `tfsdk:"callback_domain"`
	PKCE                   boolattr.Type                `tfsdk:"pkce"`
	AccessType             stringattr.Type              `tfsdk:"access_type"`
	Prompt                 strsetattr.Type              `tfsdk:"prompt"`
}

func (m *OutboundAppModel) OutboundApp(ctx context.Context, diagnostics *diag.Diagnostics) *descope.OutboundApp {
	handler := helpers.NewHandler(ctx, diagnostics)
	authorizationURLParams := make([]descope.URLParam, 0, len(m.AuthorizationURLParams.Elements()))
	for param := range listattr.Iterator[URLParamModel](m.AuthorizationURLParams, handler) {
		authorizationURLParams = append(authorizationURLParams, descope.URLParam{Key: param.Key.ValueString(), Value: param.Value.ValueString()})
	}
	tokenURLParams := make([]descope.URLParam, 0, len(m.TokenURLParams.Elements()))
	for param := range listattr.Iterator[URLParamModel](m.TokenURLParams, handler) {
		tokenURLParams = append(tokenURLParams, descope.URLParam{Key: param.Key.ValueString(), Value: param.Value.ValueString()})
	}
	defaultScopes := make([]string, 0, len(m.DefaultScopes.Elements()))
	for scope := range strsetattr.Iterator(m.DefaultScopes, handler) {
		defaultScopes = append(defaultScopes, scope)
	}
	prompt := make([]descope.PromptType, 0, len(m.Prompt.Elements()))
	for value := range strsetattr.Iterator(m.Prompt, handler) {
		prompt = append(prompt, descope.PromptType(value))
	}

	return &descope.OutboundApp{
		ID:                     m.ID.ValueString(),
		Name:                   m.Name.ValueString(),
		Description:            m.Description.ValueString(),
		ClientID:               m.ClientID.ValueString(),
		Logo:                   m.Logo.ValueString(),
		DiscoveryURL:           m.DiscoveryURL.ValueString(),
		AuthorizationURL:       m.AuthorizationURL.ValueString(),
		AuthorizationURLParams: authorizationURLParams,
		TokenURL:               m.TokenURL.ValueString(),
		TokenURLParams:         tokenURLParams,
		RevocationURL:          m.RevocationURL.ValueString(),
		DefaultScopes:          defaultScopes,
		DefaultRedirectURL:     m.DefaultRedirectURL.ValueString(),
		CallbackDomain:         m.CallbackDomain.ValueString(),
		Pkce:                   m.PKCE.ValueBool(),
		AccessType:             descope.AccessType(m.AccessType.ValueString()),
		Prompt:                 prompt,
	}
}

func (m *OutboundAppModel) SetOutboundApp(ctx context.Context, app *descope.OutboundApp) {
	m.ID = stringattr.Value(app.ID)
	m.Name = stringattr.Value(app.Name)
	m.Description = stringattr.Value(app.Description)
	m.ClientID = stringattr.Value(app.ClientID)
	m.ClientSecret = types.StringNull()
	m.Logo = stringattr.Value(app.Logo)
	m.DiscoveryURL = stringattr.Value(app.DiscoveryURL)
	m.AuthorizationURL = stringattr.Value(app.AuthorizationURL)
	m.AuthorizationURLParams = urlParamsValue(ctx, app.AuthorizationURLParams)
	m.TokenURL = stringattr.Value(app.TokenURL)
	m.TokenURLParams = urlParamsValue(ctx, app.TokenURLParams)
	m.RevocationURL = stringattr.Value(app.RevocationURL)
	m.DefaultScopes = stringSetValue(ctx, app.DefaultScopes)
	m.DefaultRedirectURL = stringattr.Value(app.DefaultRedirectURL)
	m.CallbackDomain = stringattr.Value(app.CallbackDomain)
	m.PKCE = boolattr.Value(app.Pkce)
	m.AccessType = stringattr.Value(string(app.AccessType))
	prompt := make([]string, 0, len(app.Prompt))
	for _, value := range app.Prompt {
		prompt = append(prompt, string(value))
	}
	m.Prompt = stringSetValue(ctx, prompt)
}

func urlParamsValue(ctx context.Context, params []descope.URLParam) listattr.Type[URLParamModel] {
	values := make([]*URLParamModel, 0, len(params))
	for _, param := range params {
		values = append(values, &URLParamModel{Key: stringattr.Value(param.Key), Value: stringattr.Value(param.Value)})
	}
	return helpers.Require(listtype.NewValue(ctx, values))
}

func stringSetValue(ctx context.Context, values []string) strsetattr.Type {
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.StringValue(value))
	}
	return helpers.Require(valuesettype.NewValue[types.String](ctx, elements))
}
